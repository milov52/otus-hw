package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task, len(tasks))
	stopCh := make(chan struct{})
	var errCount int32
	var wg sync.WaitGroup
	var once sync.Once

	worker := func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-taskCh:
				if !ok {
					return
				}
				if err := task(); err != nil {
					if atomic.AddInt32(&errCount, 1) >= int32(m) {
						once.Do(func() {
							close(stopCh)
						})
						return
					}
				}
			case <-stopCh:
				return
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	for _, task := range tasks {
		select {
		case <-stopCh:
			break
		case taskCh <- task:
		}
	}
	close(taskCh)
	wg.Wait()

	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
