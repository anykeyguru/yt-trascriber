package main

import (
	"log"

	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	pipeline.NewYtDl()
}
func main() {

	p := tea.NewProgram(
		tui.InitialModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
