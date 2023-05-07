package app

import (
    "fmt"
    "github.com/cihub/seelog"
    "muchat-go/chatgpt/api_base"
    "muchat-go/config"
    "muchat-go/models"
    "muchat-go/repo/censor"
    "time"
)

type JobParam struct {
    Account       *models.OpenAiAccount
    Censor        *censor.Censor
    Title         string
    TitleIndex    int
    Error         error
    InnerError    error  // 内部错误
    InnerErrorMsg string // 内部错误信息
    ErrCnt        int
    ErrCode       int
    Message       string
    UsedCap       int
    BadCnt        int
    Slug          string
    EnqueueTime   time.Time
    DoneTime      time.Time
    Messages      []api_base.ChatMessage
    PresetPrompt  *api_base.PresetPrompt
}

func (r *JobParam) GetIdStr() string {
    return fmt.Sprintf("slug=%v acc=%v", r.Slug, r.Account.Email)
}

type JobResult struct {
    JobParam *JobParam
}

type Runner struct {
    Jobs       chan *JobParam
    ResultJobs chan *JobParam
    Config     *config.Configuration
    Threads    []*RunnerThread
    //FailedThread *FailedThread
    JobStorage map[int]*JobParam
}

func NewRunner(config *config.Configuration) (r *Runner) {
    r = new(Runner)
    r.Jobs = make(chan *JobParam, config.Concurrency)
    r.ResultJobs = make(chan *JobParam, config.Concurrency*2)
    r.JobStorage = make(map[int]*JobParam)
    r.Config = config
    return
}

func (r *Runner) AddJob(titleIndex int, title string, acc *models.OpenAiAccount, censor *censor.Censor, slug string, messages []api_base.ChatMessage, presetPrompt *api_base.PresetPrompt) {
    r.Jobs <- &JobParam{
        Title:        title,
        TitleIndex:   titleIndex,
        Account:      acc,
        Censor:       censor,
        Slug:         slug,
        EnqueueTime:  time.Now(),
        Messages:     messages,
        PresetPrompt: presetPrompt,
    }
}

func (r *Runner) WaitForJob(titleIndex int) (job *JobParam) {
    tick := time.NewTicker(100 * time.Millisecond)
    dtBeg := time.Now()
    var found bool
    for {
        select {
        case <-tick.C:
            job, found = r.JobStorage[titleIndex]
            if found {
                delete(r.JobStorage, titleIndex)
                return
            }
            if dtBeg.Add(time.Second * 90).Before(time.Now()) {
                seelog.Warnf("can not wait for job: %v", titleIndex)
                job = nil
                return
            }
        }
    }
}

func (r *Runner) HandleResultJobs() {
    for {
        select {
        case job := <-r.ResultJobs:
            // 清理过旧数据
            toDeleteIndex := []int{}
            for index, oldJob := range r.JobStorage {
                deltaMinutes := time.Now().Sub(oldJob.DoneTime).Minutes()
                if deltaMinutes > 2.0 {
                    seelog.Infof("deleting too old job: %+v", *oldJob)
                    toDeleteIndex = append(toDeleteIndex, index)
                }
            }
            if len(toDeleteIndex) != 0 {
                for _, titleIndex := range toDeleteIndex {
                    delete(r.JobStorage, titleIndex)
                }
            }
            r.JobStorage[job.TitleIndex] = job
            seelog.Infof("current size of job storage: %v", len(r.JobStorage))
        }
    }
}

func (r *Runner) WaitForEmptyJob() {
    done := make(chan bool)
    go func() {
        for {
            time.Sleep(200 * time.Millisecond)
            if len(r.Jobs) == 0 && len(r.ResultJobs) == 0 {
                done <- true
                break
            }
        }
    }()
    <-done
}

func (r *Runner) Stop() {
    seelog.Infof("正在暂停 Runner ......")
    for _, t := range r.Threads {
        go t.Stop()
    }
    for _, t := range r.Threads {
        t.WaitForStop()
    }

    close(r.ResultJobs)
    //go r.FailedThread.Stop()
    //r.FailedThread.WaitForStop()
    seelog.Infof("暂停成功，所有线程已停止")
}

func (r *Runner) Start() {
    seelog.Infof("正在启动 Runner ......")

    //r.FailedThread = NewFailedThread(100, r.ResultJobs)
    //go r.FailedThread.Start()

    go r.HandleResultJobs()

    r.Threads = make([]*RunnerThread, 0, r.Config.Concurrency)
    for i := 0; i < r.Config.Concurrency; i++ {
        t := NewRunnerThread(i, r.Config, r.Jobs, r.ResultJobs)
        r.Threads = append(r.Threads, t)
        go t.Start()
    }
}
