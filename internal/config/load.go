package config

import (
	"log"

	"github.com/anykeyguru/yt-trascriber/internal/pipeline"
	"github.com/spf13/viper"
)

type AppInfo struct {
	Version string
	Commit  string
	Date    string
}

func MustLoad(info AppInfo) AppConfig {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config read error: %v", err)
	}

	var fc fileConfig
	if err := v.Unmarshal(&fc); err != nil {
		log.Fatalf("config unmarshal error: %v", err)
	}

	if fc.TextOutput == "" {
		fc.TextOutput = pipeline.DefaultOutDir
	}
	if fc.Transcriber.Language == "" {
		fc.Transcriber.Language = "auto"
	}
	if fc.Youtube.YtDlPath == "" {

	}
	return AppConfig{
		OutDir:       fc.TextOutput,
		Language:     fc.Transcriber.Language,
		YtDlpPath:    fc.Youtube.YtDlPath,
		FfmpegPath:   fc.Ffmpeg.Path,
		WhisperBin:   fc.Transcriber.Whisper.Cli,
		WhisperModel: fc.Transcriber.Whisper.Model,
		OpenAIModel:  fc.Transcriber.OpenAI.Model,
		Version:      info.Version,
		Commit:       info.Commit,
		Date:         info.Date,
	}
}
