package io_util

import (
	"bytes"
	"sync"
)

type MutexWriter struct {
	Lock     *sync.Mutex
	Buffer   *bytes.Buffer
	ChString chan string
}

func NewMutexWriter() (r *MutexWriter) {
	r = new(MutexWriter)
	r.Lock = new(sync.Mutex)
	r.Buffer = new(bytes.Buffer)
	r.ChString = make(chan string, 100)
	return
}

func (r *MutexWriter) Write(bs []byte) (n int, err error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	n, err = r.Buffer.Write(bs)
	go func() {
		r.ChString <- string(bs)
	}()
	return
}
func (r *MutexWriter) String() (text string) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	text = r.Buffer.String()
	//r.Buffer.Reset()
	return
}
