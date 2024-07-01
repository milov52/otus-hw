package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 || n == 0 {
		return nil
	}

	taskCh := make(chan Task)
	var errCount int32
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for task := range taskCh {
			if err := task(); err != nil {
				atomic.AddInt32(&errCount, 1)
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
