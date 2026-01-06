package tui

import (
	"context"
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func listenJob(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			// канал закрыт, ничего не возвращаем
			return nil
		}
		return msg
	}
}

func runPipelineJob(jobID int, url string, backend string, cfg AppConfig, ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		go func() {
			defer close(ch)

			ctx := context.Background()

			opt := pipeline.Options{
				URL:        url,
				OutDir:     cfg.OutDir,
				Language:   cfg.Language,
				Timeout:    60 * time.Minute,
				YtDlpPath:  cfg.YtDlpPath,
				FfmpegPath: cfg.FfmpegPath,

				WhisperBin:   cfg.WhisperBin,
				WhisperModel: cfg.WhisperModel,

				OpenAIKey:   cfg.OpenAIKey,
				OpenAIModel: cfg.OpenAIModel,
			}

			switch backend {
			case "whisper.cpp":
				opt.Backend = pipeline.BackendWhisperCPP
			case "openai":
				opt.Backend = pipeline.BackendOpenAI
			default:
				opt.Backend = pipeline.BackendWhisperCPP
			}

			res, err := pipeline.Run(ctx, opt, func(line string) {
				ch <- logMsg{jobID: jobID, line: line}
			})
			if err != nil {
				ch <- doneMsg{jobID: jobID, err: err}
				return
			}
			ch <- doneMsg{jobID: jobID, title: res.Title, text: res.Text}
		}()

		// стартуем “первую подписку” на канал
		return listenJob(ch)()
	}
}

func deleteWithDelayCmd(store *storage.Store, id int64) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1200 * time.Millisecond)

		err := store.Delete(id)
		return deleteDoneMsg{
			id:  id,
			err: err,
		}
	}
}
