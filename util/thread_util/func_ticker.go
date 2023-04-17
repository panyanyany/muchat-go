package thread_util

import (
    "github.com/cihub/seelog"
    "time"
)

type FuncTicker struct {
    Interval time.Duration
    Running  bool
    Handler  func()
    stopped  chan bool
}

func NewFuncTicker(interval time.Duration) (r *FuncTicker) {
    r = new(FuncTicker)
    r.Interval = interval
    r.stopped = make(chan bool)
    return
}

func (r *FuncTicker) Start() {
    r.Running = true
    go r.loop()
}

func (r *FuncTicker) Stop() {
    seelog.Infof("func ticker stopping")
    r.Running = false
    <-r.stopped
    seelog.Infof("func ticker stopped")
}

func (r *FuncTicker) loop() {
    tick := time.NewTicker(r.Interval)
    for {
        if !r.Running {
            r.stopped <- true
            break
        }
        select {
        case <-tick.C:
            r.Handler()
        }
    }
}
