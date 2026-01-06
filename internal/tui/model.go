package tui

import (
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type AppConfig struct {
	OutDir         string
	CustomPathSave string
	Language       string
	YtDlpPath      string
	FfmpegPath     string
	WhisperBin     string
	WhisperModel   string

	OpenAIKey   string
	OpenAIModel string
}

type JobStatus int

const (
	Idle JobStatus = iota
	Running
	Done
	Error
)

func (s JobStatus) String() string {
	switch s {
	case Idle:
		return "idle"
	case Running:
		return "running"
	case Done:
		return "done"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

type Job struct {
	ID        int
	URL       string
	Backend   string
	Status    JobStatus
	Logs      []string
	StartedAt time.Time
}

type ViewMode int

const (
	ViewJobs ViewMode = iota
	ViewHistory
	ViewDeleted
)

type InputMode int

const (
	InputNone InputMode = iota
	InputURL
	InputSavePath
)

type Model struct {
	InputMode InputMode
	Cfg       AppConfig

	Jobs     []Job
	Selected int
	Backend  string

	ViewMode ViewMode
	Store    *storage.Store
	History  []storage.Transcript

	Message   string
	InputPath string

	InputURL string

	jobCh map[int]chan tea.Msg
}
