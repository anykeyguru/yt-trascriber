package tui

import (
	"os"
	"path/filepath"

	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	selectedRow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	statusIdle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	statusRun   = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	statusDone  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	statusError = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func InitialModel() Model {
	store, err := storage.Open("ytranscribe.db")
	if err != nil {
		panic(err)
	}

	history, _ := store.List()

	outPath := "./out"
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		panic(err)
	}

	return Model{
		Cfg: AppConfig{
			OutDir:       outPath,
			Language:     "ru",
			YtDlpPath:    "yt-dlp",
			FfmpegPath:   "ffmpeg",
			WhisperBin:   "./ext-tools/whisper/whisper-cli",
			WhisperModel: "./models/ggml-medium.bin",
			OpenAIModel:  "gpt-4o-transcribe",
		},
		Backend:  "whisper.cpp",
		Store:    store,
		History:  history,
		ViewMode: ViewJobs,
		jobCh:    make(map[int]chan tea.Msg),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
