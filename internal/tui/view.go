package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	switch m.ViewMode {
	case ViewJobs:
		return m.jobsView()
	case ViewHistory:
		return m.historyView()
	case ViewDeleted:
		return m.deletedRecordView()
	default:
		return m.jobsView()
	}
}

func (m Model) jobsView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ytranscribe — TUI\n\n"))
	b.WriteString(fmt.Sprintf(
		"Backend: %s   (press 'b' to switch)\n\n",
		m.Backend,
	))

	b.WriteString("Jobs:\n")

	for i, j := range m.Jobs {
		row := fmt.Sprintf(
			"[%d] %-8s %-10s %s",
			j.ID,
			j.Status.String(),
			j.Backend,
			j.URL,
		)

		if i == m.Selected {
			row = selectedRow.Render("> " + row)
		} else {
			switch j.Status {
			case Idle:
				row = statusIdle.Render("  " + row)
			case Running:
				row = statusRun.Render("  " + row)
			case Done:
				row = statusDone.Render("  " + row)
			case Error:
				row = statusError.Render("  " + row)
			default:
				row = "  " + row
			}

		}
		b.WriteString(row + "\n")
	}

	b.WriteString("\nLogs:\n")

	if len(m.Jobs) > 0 && m.Selected < len(m.Jobs) {
		for _, l := range m.Jobs[m.Selected].Logs {
			b.WriteString("  " + l + "\n")
		}
	}

	if m.Message != "" {
		b.WriteString("\n" + m.Message + "\n")
		b.WriteString("> " + m.InputURL)
	}

	b.WriteString("\n\nKeys: a-add  enter-run  b-backend  h-history  q-quit")

	return b.String()
}

func (m Model) historyView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("History (SQLite)\n\n"))

	//var errList error
	//m.History, errList = m.Store.List()
	//if errList != nil {
	//	b.WriteString(statusError.Render(errList.Error()))
	//	return b.String()
	//}
	b.WriteString(fmt.Sprintf("Total: %d\n\n", len(m.History)))
	b.WriteString("List:\n")
	for i, t := range m.History {
		row := fmt.Sprintf(
			"[%d] %-10s %s",
			t.ID,
			t.Backend,
			t.Title,
		)

		if i == m.Selected {
			row = statusDone.Render("> " + row)
		} else {
			row = "  " + row
		}
		b.WriteString(row + "\n")
	}

	if m.InputMode == InputSavePath {
		b.WriteString("\n" + m.Message + "\n")
		b.WriteString("> " + m.InputPath + "\n")
	}

	b.WriteString("\n\nKeys: h-back  s-save to file  d-delete record  q-quit")

	return b.String()
}

func (m Model) deletedRecordView() string {
	return titleStyle.Render("Deleting…\n\n") +
		m.Message +
		statusError.Render("\n\n(deleting...)")
}
