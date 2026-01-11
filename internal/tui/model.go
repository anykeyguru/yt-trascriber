package tui

import (
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/config"
	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

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
	Backend   pipeline.Backend
	Status    JobStatus
	Logs      []string
	StartedAt time.Time
}

type ViewMode int

const (
	ViewJobs ViewMode = iota
	ViewHistory
	ViewDeleted
	ViewMP3
)

type InputMode int

const (
	InputNone InputMode = iota
	InputURL
	InputSavePath
)

type Model struct {
	InputMode InputMode
	Cfg       config.AppConfig

	Jobs     []Job
	Selected int
	Backend  pipeline.Backend

	ViewMode ViewMode
	Store    *storage.Store
	History  []storage.Transcript

	Message   string
	InputPath string

	InputURL string

	jobCh map[int]chan tea.Msg

	Width  int
	Height int
}
