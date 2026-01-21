package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/helpers"
	"github.com/google/uuid"
)

func RunGrubMP3(ctx context.Context, opt Options, log LogFn) (*Result, error) {
	opt.URL = normalizeURL(opt.URL)

	if opt.URL == "" {
		return nil, errors.New("empty URL")
	}

	log(fmt.Sprintf("URL: %s", opt.URL))
	if res, bad := helpers.ValidateURL(opt.URL); bad {
		return nil, errors.New(res)
	}

	if opt.Timeout <= 0 {
		opt.Timeout = 60 * time.Minute
	}
	if opt.YtDlpPath == "" {
		opt.YtDlpPath = defaultYtDlpPath
	}
	if opt.FfmpegPath == "" {
		opt.FfmpegPath = defaultFfmpegPath
	}
	if opt.OutDir == "" {
		opt.OutDir = DefaultOutDir
	}

	ctx, cancel := context.WithTimeout(ctx, opt.Timeout)
	defer cancel()

	if err := os.MkdirAll(opt.OutDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(BaseTempDir, 0o755); err != nil {
		return nil, err
	}
	
	tmpDir, err := os.MkdirTemp(BaseTempDir, "ytranscribe-*")
	if err != nil {
		return nil, err
	}
	if !opt.KeepTemp {
		defer os.RemoveAll(tmpDir)
	} else {
		log(fmt.Sprintf("keep temp dir: %s", tmpDir))
	}

	title, err := YtFetcher.GetTitle(ctx, opt, log)
	//title, err := ytTitle(ctx, opt, log)
	if err != nil {
		return nil, err
	}
	if title == "" {
		title = "youtube"
	}
	log("title: " + title)

	audioPath, err := YtFetcher.FetchAudio(ctx, opt, tmpDir, log)
	if err != nil {
		return nil, err
	}
	log("audio downloaded: " + audioPath)

	mp3Path := filepath.Join("./audio", "audio"+uuid.New().String()[:8]+".mp3")
	if err := ffmpegToMP3Stereo(ctx, opt, audioPath, mp3Path, log); err != nil {
		return nil, err
	}
	log("audio converted: " + mp3Path)

	err = os.Rename(mp3Path, filepath.Join("./audio", title+".mp3"))
	if err != nil {
		return nil, err
	}

	return &Result{Title: title, Text: ""}, nil
}
