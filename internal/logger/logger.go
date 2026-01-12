package logger

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Inst *log.Logger

func New(opts Options) (*log.Logger, func(), error) {
	if opts.File == "" {
		return NewNoop(), func() {}, nil
	}

	lj := &lumberjack.Logger{
		Filename:   opts.File,
		MaxSize:    defaultInt(opts.MaxSizeMB, 50),
		MaxBackups: defaultInt(opts.MaxBackups, 5),
		MaxAge:     defaultInt(opts.MaxAgeDays, 28),
		Compress:   opts.Compress,
	}

	var w io.Writer = lj
	if opts.AlsoStdout {
		w = io.MultiWriter(os.Stdout, lj)
	}

	logger := log.New(
		w,
		opts.LevelPrefix,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile,
	)

	cleanup := func() {
		_ = lj.Close()
	}

	Inst = logger
	return logger, cleanup, nil
}

func defaultInt(v, d int) int {
	if v > 0 {
		return v
	}
	return d
}
