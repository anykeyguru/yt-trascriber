package pipeline

import "context"

var whisperQueue = make(chan whisperTask, 16)

type JobStatus int

type whisperTask struct {
	ctx   context.Context
	opt   Options
	wav   string
	title string
	log   LogFn

	result chan whisperResult
}
type whisperResult struct {
	text string
	err  error
}

func StartWhisperWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case task := <-whisperQueue:
				text, err := whisperCPP(
					task.ctx,
					task.opt,
					task.wav,
					task.title,
					task.log,
				)

				task.result <- whisperResult{
					text: text,
					err:  err,
				}
			}
		}
	}()
}
