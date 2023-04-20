package app

import (
    "github.com/cihub/seelog"
    "muchat-go/util/file_util"
)

type FailedThread struct {
    FailedJobs chan *JobParam
    ShouldStop chan bool
    Stopped    chan bool
    Running    bool
    Index      int
}

func NewFailedThread(index int, failedJobs chan *JobParam) (r *FailedThread) {
    r = new(FailedThread)
    r.Index = index
    r.FailedJobs = failedJobs
    r.ShouldStop = make(chan bool)
    r.Stopped = make(chan bool)
    return
}

func (r *FailedThread) Stop() {
    seelog.Infof("正在暂停线程 #%v......", r.Index)
    r.ShouldStop <- true
}
func (r *FailedThread) WaitForStop() {
    <-r.Stopped
    seelog.Infof("已退出线程 #%v......", r.Index)
}

func (r *FailedThread) Start() {
    for {
        select {
        case <-r.ShouldStop:
            for job := range r.FailedJobs {
                seelog.Infof("job:%v", job)
                r.DoJob(job)
            }
            r.Stopped <- true
            return
        case job := <-r.FailedJobs:
            // 这里会产生 job=nil ？
            r.DoJob(job)
        }
    }
}

func (r *FailedThread) DoJob(job *JobParam) {
    if job == nil {
        return
    }
    seelog.Debugf("failed threads：len=%v, cap=%v", len(r.FailedJobs), cap(r.FailedJobs))
    seelog.Infof("保存失败题目：%v", job.Title)
    file_util.AppendFile("失败题目.txt", job.Title+"\n", nil)
}
