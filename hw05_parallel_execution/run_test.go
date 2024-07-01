package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("zero workers", func(t *testing.T) {
		tasks := []Task{
			func() error { return nil },
			func() error { return errors.New("error") },
		}

		err := Run(tasks, 0, 1)
		require.NoError(t, err, "expected no error when no workers")
	})

	t.Run("zero maximum errors allowed", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("error") },
			func() error { return nil },
		}

		err := Run(tasks, 2, 0)
		require.True(t, errors.Is(err, ErrErrorsLimitExceeded), "expected ErrErrorsLimitExceeded when max errors allowed is 0")
	})

	t.Run("no tasks", func(t *testing.T) {
		var tasks []Task

		err := Run(tasks, 5, 1)
		require.NoError(t, err, "expected no error when no tasks")
	})

	t.Run("less tasks than workers", func(t *testing.T) {
		var runTasksCount int32

		tasks := []Task{
			func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			},
			func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			},
		}

		err := Run(tasks, 5, 1)
		require.NoError(t, err)
		require.Equal(t, int32(2), runTasksCount, "all tasks should be completed")
	})

	t.Run("long running tasks", func(t *testing.T) {
		tasks := []Task{
			func() error {
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			func() error {
				time.Sleep(300 * time.Millisecond)
				return nil
			},
		}

		start := time.Now()
		err := Run(tasks, 2, 1)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.GreaterOrEqual(t, elapsedTime, 300*time.Millisecond, "tasks should have taken at least 300ms to complete")
	})

	t.Run("all tasks succeeding quickly", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks[i] = func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			}
		}

		err := Run(tasks, 5, 1)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "all tasks should be completed")
	})

}
