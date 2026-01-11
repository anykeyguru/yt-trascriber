package main

import (
	"context"
	"log"

	"github.com/anykeyguru/yt-trascriber/internal/config"
	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	pipeline.NewYtDl()
}
func main() {

	cfg := config.AppConfig{
		OutDir:       pipeline.DefaultOutDir,
		Language:     "ru",
		YtDlpPath:    "yt-dlp",
		FfmpegPath:   "ffmpeg",
		WhisperBin:   "./ext-tools/whisper/whisper-cli",
		WhisperModel: "./models/ggml-medium.bin",
		OpenAIModel:  "gpt-4o-transcribe",
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline.StartWhisperWorker(ctx)

	p := tea.NewProgram(
		tui.InitialModel(cfg),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
