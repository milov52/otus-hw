package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func executeStage(done In, in In, stage Stage) Out {
	out := make(Bi)
	stageOut := stage(in)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				for range stageOut {
				}
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

func workTask(done In, in In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = executeStage(done, out, stage)
	}
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := workTask(done, in, stages...)

	result := make(Bi)
	go func() {
		defer close(result)
		for {
			select {
			case <-done:
				return
			case v, ok := <-out:
				if !ok {
					return
				}
				result <- v
			}
		}
	}()
	return result
}
