package pipeline

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/anykeyguru/yt-trascriber/internal/helpers"
)

type Backend string

const (
	BackendWhisperCPP Backend = "whisper.cpp"
	BackendOpenAI     Backend = "openai"
	BackendFFMPEG     Backend = "ffmpeg"
	defaultYtDlpPath          = "yt-dlp"
	defaultFfmpegPath         = "ffmpeg"
	DefaultOutDir             = "./out"
)

func (b Backend) String() string {
	return string(b)
}

type Options struct {
	URL      string
	Backend  Backend
	Language string // "auto", "ru", "en", ...
	OutDir   string // where to store exported results
	KeepTemp bool
	Timeout  time.Duration

	YtDlpPath  string
	FfmpegPath string

	// whisper.cpp
	WhisperBin   string // whisper-cli
	WhisperModel string

	// openai (будет подключено следующим шагом)
	OpenAIKey   string
	OpenAIModel string
}

type Result struct {
	Title string
	Text  string
}

type LogFn func(line string)

func RunGrubTExtFromYTb(ctx context.Context, opt Options, log LogFn) (*Result, error) {
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

	tmpDir, err := os.MkdirTemp("", "ytranscribe-*")
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

	//audioPath, err := ytDownloadAudio(ctx, opt, tmpDir, log)
	audioPath, err := YtFetcher.FetchAudio(ctx, opt, tmpDir, log)
	if err != nil {
		return nil, err
	}
	log("audio downloaded: " + audioPath)

	wavPath := filepath.Join(tmpDir, "audio.wav")
	if err := ffmpegToWav16kMono(ctx, opt, audioPath, wavPath, log); err != nil {
		return nil, err
	}
	log("audio converted: " + wavPath)

	if err := os.WriteFile("audio.log", []byte(wavPath), 0644); err != nil {
		return nil, err
	}

	switch opt.Backend {
	case BackendWhisperCPP:
		resCh := make(chan whisperResult, 1)

		whisperQueue <- whisperTask{
			ctx:    ctx,
			opt:    opt,
			wav:    wavPath,
			title:  title,
			log:    log,
			result: resCh,
		}

		res := <-resCh
		if res.err != nil {
			return nil, res.err
		}
		return &Result{Title: title, Text: res.text}, nil

	case BackendOpenAI:
		return nil, errors.New("openai backend: not implemented in this step (next step we’ll add API call)")

	default:
		return nil, fmt.Errorf("unknown backend: %s", opt.Backend)
	}
}

func streamCmd(cmd *exec.Cmd, log LogFn) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan struct{}, 2)
	go func() {
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			log(sc.Text())
		}
		done <- struct{}{}
	}()
	go func() {
		sc := bufio.NewScanner(stderr)
		for sc.Scan() {
			log(sc.Text())
		}
		done <- struct{}{}
	}()

	<-done
	<-done

	return cmd.Wait()
}
