package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	errCh := make(chan error)
	errThresholdReached := make(chan struct{})
	tasksCompleted := make(chan struct{})

	// Горутина для контроля ошибок
	go func() {
		defer close(errThresholdReached)
		errCount := 0
		for {
			select {
			case <-tasksCompleted:
				return
			case err := <-errCh:
				if err != nil {
					errCount++
					if errCount >= m {
						return
					}
				}
			}
		}
	}()

	// Запуск воркеров
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-errThresholdReached:
					return
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					if err := task(); err != nil {
						select {
						case errCh <- err:
						case <-errThresholdReached:
							return
						}
					} else {
					}
				}
			}
		}()
	}

	// Добавление задач
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case <-errThresholdReached:
				return
			case taskCh <- task:
			}
		}
	}()

	wg.Wait()
	close(tasksCompleted)

	select {
	case <-errThresholdReached:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
