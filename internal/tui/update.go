package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/helpers"
	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type logMsg struct {
	jobID int
	line  string
}

type doneMsg struct {
	jobID int
	title string
	text  string
	err   error
}

type deleteDoneMsg struct {
	id  int64
	err error
}

// Update processes incoming messages to modify the state of the model and potentially produce new commands for execution.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.InputMode == InputURL {

			switch msg.Type {

			case tea.KeyEsc:
				m.InputMode = InputNone
				m.InputURL = ""
				m.InputPath = ""
				m.Message = ""
				return m, nil

			case tea.KeyRunes:
				m.InputURL += string(msg.Runes)

			case tea.KeyBackspace:
				if len(m.InputURL) > 0 {
					m.InputURL = m.InputURL[:len(m.InputURL)-1]
				}

			case tea.KeyEnter:
				url := strings.TrimSpace(m.InputURL)

				if url == "" {
					m.Message = "URL is empty"
					return m, nil
				}

				if res, bad := helpers.ValidateURL(url); bad {
					m.Message = res
					return m, nil
				}

				job := Job{
					ID:      len(m.Jobs) + 1,
					URL:     url,
					Backend: m.Backend,
					Status:  Idle,
				}

				m.Jobs = append(m.Jobs, job)

				m.InputURL = ""
				m.InputMode = InputNone
				m.Message = ""
			}

			return m, nil
		}

		if m.InputMode == InputSavePath {
			switch msg.Type {
			case tea.KeyRunes:
				m.InputPath += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.InputPath) > 0 {
					m.InputPath = m.InputPath[:len(m.InputPath)-1]
				}
			case tea.KeyEnter:
				// export to file
				if strings.TrimSpace(m.InputPath) == "" {
					m.Message = "Path is empty"
					return m, nil
				}
				t := m.History[m.Selected]

				if err := os.MkdirAll(filepath.Dir(m.InputPath), 0755); err != nil {
					m.Message = "Create dir failed: " + err.Error()
					return m, nil
				}

				err := os.WriteFile(m.InputPath, []byte(t.Text), 0644)
				if err != nil {
					m.Message = "Save failed: " + err.Error()
				} else {
					m.Message = "Saved to " + m.InputPath
				}
				path := m.InputPath

				m.InputPath = ""
				m.InputMode = InputNone
				m.Message = "saved to: " + path

			}
			return m, nil
		}

		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up":
			switch m.ViewMode {
			case ViewJobs:
				if m.Selected > 0 {
					m.Selected--
				}
			case ViewHistory:
				if m.Selected > 0 {
					m.Selected--
				}
			}

		case "down":
			switch m.ViewMode {
			case ViewJobs:
				if m.Selected < len(m.Jobs)-1 {
					m.Selected++
				}
			case ViewHistory:
				if m.Selected < len(m.History)-1 {
					m.Selected++
				}
			}

		case "b":
			if m.ViewMode == ViewJobs && m.InputMode == InputNone {
				if m.Backend == "whisper.cpp" {
					m.Backend = "openai"
				} else {
					m.Backend = "whisper.cpp"
				}
			}

		case "a":
			if m.ViewMode == ViewJobs && m.InputMode == InputNone {
				m.InputURL = ""
				m.InputMode = InputURL
				m.Message = "Enter YouTube URL (press Enter)"
			}

		case "h":
			m.InputMode = InputNone
			m.InputPath = ""
			m.InputURL = ""
			m.Message = ""

			if m.ViewMode == ViewJobs {
				m.ViewMode = ViewHistory
				if m.Selected >= len(m.History) {
					m.Selected = max(len(m.History)-1, 0)
				}
			} else {
				m.ViewMode = ViewJobs
				if m.Selected >= len(m.Jobs) {
					m.Selected = max(len(m.Jobs)-1, 0)
				}
			}

		case "s":
			if m.ViewMode == ViewHistory && len(m.History) > 0 {
				t := m.History[m.Selected]

				m.InputPath = filepath.Join(
					m.Cfg.OutDir,
					pipeline.CreateFileNameForSave(t.Title)+".txt",
				)

				m.InputMode = InputSavePath
				m.Message = "Save to file: press Enter or edit path"
			}

		case "d":
			if m.ViewMode == ViewHistory && len(m.History) > 0 {
				t := m.History[m.Selected]

				m.ViewMode = ViewDeleted
				m.Message = fmt.Sprintf("Record %d will be deleted…", t.ID)

				return m, deleteWithDelayCmd(m.Store, t.ID)
			}

		case "enter":
			if m.ViewMode == ViewJobs && len(m.Jobs) > 0 {
				job := &m.Jobs[m.Selected]

				if job.Status == Idle {
					job.Status = Running

					ch := make(chan tea.Msg, 200)
					m.jobCh[job.ID] = ch

					return m, runPipelineJob(
						job.ID,
						job.URL,
						job.Backend,
						m.Cfg,
						ch,
					)
				}
			}

		}

	case deleteDoneMsg:
		if msg.err != nil {
			m.Message = "Delete failed: " + msg.err.Error()
		} else {
			m.Message = fmt.Sprintf("Record %d deleted", msg.id)
		}

		m.ViewMode = ViewHistory
		m.History, _ = m.Store.List()

		if m.Selected >= len(m.History) {
			m.Selected = max(len(m.History)-1, 0)
		}

		return m, nil

	case logMsg:
		for i := range m.Jobs {
			if m.Jobs[i].ID == msg.jobID {
				m.Jobs[i].Logs = append(m.Jobs[i].Logs, msg.line)
				// удерживаем последние N строк, чтобы память не росла
				if len(m.Jobs[i].Logs) > 300 {
					m.Jobs[i].Logs = m.Jobs[i].Logs[len(m.Jobs[i].Logs)-300:]
				}
				break
			}
		}
		// продолжить слушать этот job channel
		ch := m.jobCh[msg.jobID]
		if ch != nil {
			return m, listenJob(ch)
		}
		return m, nil

	case doneMsg:
		for i := range m.Jobs {
			if m.Jobs[i].ID == msg.jobID {
				if msg.err != nil {
					m.Jobs[i].Status = Error
					m.Jobs[i].Logs = append(m.Jobs[i].Logs, "ERROR: "+msg.err.Error())
				} else {
					m.Jobs[i].Status = Done
					m.Jobs[i].Logs = append(m.Jobs[i].Logs, "done: "+msg.title)
					// пока просто положим текст в logs хвостом (или добавь поле Job.Text)
					if strings.TrimSpace(msg.text) != "" {
						m.Jobs[i].Logs = append(m.Jobs[i].Logs, "----- transcript (first 20 lines) -----")
						lines := strings.Split(msg.text, "\n")
						for k := 0; k < len(lines) && k < 20; k++ {
							m.Jobs[i].Logs = append(m.Jobs[i].Logs, lines[k])
						}
					}
				}
				if err := m.Store.Save(storage.Transcript{
					Title:   msg.title + fmt.Sprintf("_%d", time.Now().UnixNano()),
					URL:     m.Jobs[i].URL,
					Backend: m.Jobs[i].Backend,
					Text:    msg.text,
				}); err != nil {
					m.Message = "Save failed: " + err.Error()
				}
				m.History, _ = m.Store.List()
				break
			}
		}
		// job завершился — отписываемся
		delete(m.jobCh, msg.jobID)
		return m, nil

	}

	return m, nil
}
