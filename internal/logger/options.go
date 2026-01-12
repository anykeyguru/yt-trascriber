package logger

type Options struct {
	File        string
	LevelPrefix string

	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool

	AlsoStdout bool
}
