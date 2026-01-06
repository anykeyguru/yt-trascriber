package pipeline

import (
	"context"
)

type YouTubeFetcher interface {
	GetTitle(ctx context.Context, opt Options, log LogFn) (string, error)
	FetchAudio(ctx context.Context, opt Options, tmpDir string, log LogFn) (string, error)
}

var YtFetcher YouTubeFetcher
