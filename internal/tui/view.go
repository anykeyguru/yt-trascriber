package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.ViewMode {
	case ViewJobs:
		return m.jobsView()
	case ViewHistory:
		return m.historyView()
	case ViewDeleted:
		return m.deletedRecordView()
	case ViewMP3:
		return m.jobsMP3View()
	default:
		return m.jobsView()
	}
}

func (m Model) jobsView() string {

	var b strings.Builder

	b.WriteString(titleStyle.Render("ytranscribe — Jobs\n\n"))
	b.WriteString(fmt.Sprintf(
		"Backend: %s   (b = switch)   Jobs: %d\n\n",
		m.Backend,
		len(m.Jobs),
	))
	b.WriteString(m.divider())

	// ---- Header
	b.WriteString(
		headerStyle.Render(
			fmt.Sprintf(
				"  %-4s %-7s   %-12s  %s",
				"ID",
				"STATUS",
				"BACKEND",
				"URL",
			),
		),
	)
	b.WriteString(m.divider())

	b.WriteString("\n")
	// ---- Rows
	for i, j := range m.Jobs {
		prefix := "  "
		rowStyle := lipgloss.NewStyle()

		if i == m.Selected {
			prefix = "> "
			rowStyle = selectedStyle
		}

		line := fmt.Sprintf(
			"%-4d %-7s   %-14s %s",
			j.ID,
			renderStatus(j.Status),
			j.Backend,
			j.URL,
		)

		b.WriteString(rowStyle.Render(prefix + line))
		b.WriteString("\n")
	}

	// ---- Logs
	b.WriteString(m.divider())
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("Logs:"))
	b.WriteString("\n")

	if len(m.Jobs) > 0 && m.Selected < len(m.Jobs) {
		for _, l := range m.Jobs[m.Selected].Logs {
			b.WriteString("  " + l + "\n")
		}
	}

	// ---- Input / message
	if m.Message != "" {
		b.WriteString("\n" + m.Message + "\n")
		b.WriteString("> " + m.InputURL)
	}

	b.WriteString("\n\nKeys: a-add  enter-run  b-backend m-grab mp3  h-history  q-quit")

	return b.String()
}

func (m Model) historyView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("History (SQLite)\n\n"))

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

func (m Model) divider() string {
	if m.Width <= 0 {
		return "\n"
	}
	return "\n" + strings.Repeat("─", m.Width) + "\n"
}

func (m Model) jobsMP3View() string {

	var b strings.Builder

	b.WriteString(titleStyle.Render("ytranscribe — Jobs\n\n"))
	b.WriteString("MP3 Grubber")

	b.WriteString(m.divider())

	// ---- Header
	b.WriteString(
		headerStyle.Render(
			fmt.Sprintf(
				"  %-4s %-7s   %-12s  %s",
				"ID",
				"STATUS",
				"BACKEND",
				"URL",
			),
		),
	)
	b.WriteString(m.divider())

	b.WriteString("\n")
	// ---- Rows
	for i, j := range m.Jobs {
		prefix := "  "
		rowStyle := lipgloss.NewStyle()

		if i == m.Selected {
			prefix = "> "
			rowStyle = selectedStyle
		}

		line := fmt.Sprintf(
			"%-4d %-7s   %-14s %s",
			j.ID,
			renderStatus(j.Status),
			"ffmpeg",
			j.URL,
		)

		b.WriteString(rowStyle.Render(prefix + line))
		b.WriteString("\n")
	}

	// ---- Logs
	b.WriteString(m.divider())
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("Logs:"))
	b.WriteString("\n")

	if len(m.Jobs) > 0 && m.Selected < len(m.Jobs) {
		for _, l := range m.Jobs[m.Selected].Logs {
			b.WriteString("  " + l + "\n")
		}
	}

	// ---- Input / message
	if m.Message != "" {
		b.WriteString("\n" + m.Message + "\n")
		b.WriteString("> " + m.InputURL)
	}

	b.WriteString("\n\nKeys: m-back a-add  enter-run  q-quit")

	return b.String()
}
