package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anykeyguru/yt-trascriber/internal/config"
	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/anykeyguru/yt-trascriber/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	AppName = "ytranscribe"
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func init() {
	pipeline.NewYtDl()
}
func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf(
			"[%s] %s\ncommit: %s\nbuild:  %s\n",
			AppName, Version, Commit, Date,
		)
		os.Exit(0)
	}

	cfg := config.MustLoad(config.AppInfo{Version: Version, Commit: Commit, Date: Date})
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
