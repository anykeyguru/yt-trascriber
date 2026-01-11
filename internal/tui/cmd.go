package tui

import (
	"context"
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/config"
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

func runPipelineJob(jobID int, url string, backend pipeline.Backend, cfg config.AppConfig, ch chan tea.Msg) tea.Cmd {
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

				KeepTemp: false,

				WhisperBin:   cfg.WhisperBin,
				WhisperModel: cfg.WhisperModel,

				OpenAIKey:   cfg.OpenAIKey,
				OpenAIModel: cfg.OpenAIModel,
			}

			var runFn func(context.Context, pipeline.Options, pipeline.LogFn) (*pipeline.Result, error)
			switch backend {
			case "whisper.cpp":
				opt.Backend = pipeline.BackendWhisperCPP
				runFn = pipeline.RunGrubTExtFromYTb
			case "openai":
				opt.Backend = pipeline.BackendOpenAI
				runFn = pipeline.RunGrubTExtFromYTb
			case "ffmpeg":
				opt.Backend = pipeline.BackendFFMPEG
				runFn = pipeline.RunGrubMP3
			default:
				opt.Backend = pipeline.BackendWhisperCPP
				runFn = pipeline.RunGrubTExtFromYTb
			}

			res, err := runFn(ctx, opt, func(line string) {
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
