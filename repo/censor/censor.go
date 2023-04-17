package censor

import (
	"fmt"
	"github.com/cihub/seelog"
	"sync"
	"time"
)

type Censor struct {
	Adapter        Adapter
	Jobs           chan *Job
	Stopping       chan bool
	NThreads       int
	Running        bool
	JobStorage     map[int]*Job
	JobStorageLock *sync.Mutex
}

type Job struct {
	Id             int
	Text           string
	Done           bool
	Safe           bool
	Err            error
	AuditingResult *TextAuditingResult
}

func NewCensor(adapter Adapter, nThreads int) (r *Censor) {
	r = new(Censor)
	r.Adapter = adapter
	r.Jobs = make(chan *Job, nThreads)
	r.Stopping = make(chan bool, nThreads)
	r.NThreads = nThreads
	r.JobStorage = make(map[int]*Job)
	r.JobStorageLock = new(sync.Mutex)
	return
}

func (r *Censor) Start() {
	r.Running = true
	for i := 0; i < r.NThreads; i++ {
		go r.loop()
	}
}

func (r *Censor) Stop() {
	seelog.Infof("stopping censor threads")
	r.Running = false
	for i := range r.Stopping {
		_ = i
	}
	seelog.Infof("all censor threads stopped")
}

func (r *Censor) loop() {
	defer func() {
		r.Stopping <- true
	}()
	for {
		if !r.Running {
			break
		}
		select {
		case job := <-r.Jobs:
			//seelog.Infof("拿到审核任务")
			if !r.Running {
				break
			}
			res, err := r.Adapter.MakeTextAuditing(fmt.Sprintf("%v", job.Id), job.Text)

			job.Done = true
			job.Safe = res.Safe
			job.AuditingResult = res
			job.Err = err

			r.JobStorageLock.Lock()
			r.JobStorage[job.Id] = job
			r.JobStorageLock.Unlock()
		}
	}
}

func (r *Censor) WaitForJob(id int) *Job {
	tick := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			job, found := r.JobStorage[id]
			if found {
				delete(r.JobStorage, id)
				return job
			}
		}
	}
}
func (r *Censor) AddJob(id int, text string) {
	go func() {
		r.Jobs <- &Job{
			Id:   id,
			Text: text,
		}
	}()
}

func (r *Censor) CheckText(id int, text string) (safe bool, err error) {
	var res *TextAuditingResult
	res, err = r.Adapter.MakeTextAuditing(fmt.Sprintf("%v", id), text)
	if err != nil {
		return
	}
	safe = res.Safe
	return
}
