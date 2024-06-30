package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// worker запускает задания из taskCh и учитывает ошибки в счетчике errCount.
// Закрывает канал errThresholdReached, если количество ошибок превышает порог m.
func worker(
	taskCh <-chan Task, // Входной канал с задачами
	wg *sync.WaitGroup, // Группа ожидания для синхронизации горутин
	errThresholdReached chan struct{}, // Канал для оповещения о превышении порога ошибок
	once *sync.Once, // Для однократного выполнения закрытия канала
	m int, // Порог допустимых ошибок
	errCount *int32, // Счетчик ошибок

) {
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
				if atomic.AddInt32(errCount, 1) >= int32(m) {
					once.Do(func() {
						close(errThresholdReached)
					})
					return
				}
			}
		}
	}
}

// addTasks отправляет задания в taskCh или прекращает добавление, если errThresholdReached закрыт.
func addTask(taskCh chan<- Task, tasks []Task, errThresholdReached chan struct{}) {
	defer close(taskCh)
	for _, task := range tasks {
		select {
		case <-errThresholdReached:
			return
		case taskCh <- task:
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	errThresholdReached := make(chan struct{})
	tasksCompleted := make(chan struct{})
	var errCount int32
	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskCh, &wg, errThresholdReached, &once, m, &errCount)
	}

	go addTask(taskCh, tasks, errThresholdReached)
	wg.Wait()
	close(tasksCompleted)

	select {
	case <-errThresholdReached:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
