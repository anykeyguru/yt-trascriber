package logger

import (
	"io"
	"log"
)

func NewNoop() *log.Logger {
	return log.New(io.Discard, "", 0)
}
