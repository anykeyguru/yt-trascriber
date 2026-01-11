package config

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
