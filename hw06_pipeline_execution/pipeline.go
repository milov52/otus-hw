package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = <-chan interface{}
	Bi  = chan interface{}
)

type Stage func(in In) Out

func single(v interface{}) In {
	out := make(Bi, 1)
	out <- v
	close(out)
	return out
}

func workTask(done In, in In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		stageOut := make(Bi)
		go func(in In, done In, stage Stage, stageOut Bi) {
			defer close(stageOut)
			for v := range in {
				for result := range stage(single(v)) {
					select {
					case <-done:
						return
					default:
						stageOut <- result
					}
				}
			}
		}(out, done, stage, stageOut)
		out = stageOut
	}
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return workTask(done, in, stages...)
}
