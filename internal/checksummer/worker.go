/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package checksummer

import (
	"runtime"
	"samedup/internal/logger"
	"samedup/internal/node"
	"sync"
)

// Pool worker pool
type Pool struct {
	jobChan chan *job
	wg      sync.WaitGroup
	flogger *logger.FileLogger
}

type job struct {
	hashType  int
	path      string
	chksumPtr *string
	nb        int
}

func (p *Pool) worker(jobChan <-chan *job) {
	defer p.wg.Done()
	for job := range jobChan {
		chk, err := checksum(job.path, job.hashType, int64(job.nb))
		if err != nil {
			p.flogger.Error(err)
			continue
		}
		*job.chksumPtr = chk
	}
}

// AddJob push a new job to the queue
func (p *Pool) AddJob(entry *node.Node, hashType int, nb int) {
	storage := &entry.Checksum
	if len(*storage) > 0 {
		// already done
		return
	}
	j := &job{
		hashType:  hashType,
		path:      entry.AbsPath,
		chksumPtr: storage,
		nb:        nb,
	}
	p.jobChan <- j
}

// Wait wait for all worker to be done
func (p *Pool) Wait() {
	close(p.jobChan)
	p.wg.Wait()
}

// Start start the pool
func (p *Pool) Start() {
	for w := 1; w <= runtime.NumCPU(); w++ {
		p.wg.Add(1)
		go p.worker(p.jobChan)
	}
}

// NewPool create a new pool
func NewPool(flogger *logger.FileLogger) *Pool {
	p := &Pool{
		jobChan: make(chan *job),
		flogger: flogger,
	}

	return p
}
