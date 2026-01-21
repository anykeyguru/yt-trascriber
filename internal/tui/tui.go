package tui

import (
	"os"
	"path/filepath"

	"github.com/anykeyguru/yt-trascriber/internal/config"
	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	idStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	urlStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedStyle = lipgloss.NewStyle().Bold(true)

	titleStyle  = lipgloss.NewStyle().Bold(true)
	selectedRow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	statusIdle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	statusRun   = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	statusDone  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	statusError = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

func InitialModel(cfg config.AppConfig) Model {
	store, err := storage.Open("ytranscribe.db")
	if err != nil {
		panic(err)
	}

	history, _ := store.List()

	if err := os.MkdirAll(filepath.Dir(cfg.OutDir), 0755); err != nil {
		panic(err)
	}

	return Model{
		Cfg:      cfg,
		Backend:  pipeline.BackendWhisperCPP,
		Store:    store,
		History:  history,
		ViewMode: ViewMP3,
		jobCh:    make(map[int]chan tea.Msg),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func renderStatus(s JobStatus) string {
	switch s {
	case Idle:
		return statusIdle.Render(" IDLE ")
	case Running:
		return statusRun.Render(" RUN ")
	case Done:
		return statusDone.Render(" DONE ")
	case Error:
		return statusError.Render(" ERR ")
	default:
		return " ??? "
	}
}
