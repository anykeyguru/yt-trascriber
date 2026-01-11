package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func NewYtDl() {
	YtFetcher = &YtDl{}
}

type YtDl struct{}

func (yt *YtDl) GetTitle(ctx context.Context, opt Options, log LogFn) (string, error) {
	cmd := exec.CommandContext(ctx, opt.YtDlpPath, "--no-playlist", "--print", "title", opt.URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log(string(out))
		return "", fmt.Errorf("yt-dlp title: %w", err)
	}

	return Sanitize(string(out)), nil
}

func (yt *YtDl) FetchAudio(ctx context.Context, opt Options, tmpDir string, log LogFn) (string, error) {
	// Пишем в фиксированный шаблон, чтобы не зависеть от title в имени файла
	outTmpl := filepath.Join(tmpDir, "audio.%(ext)s")

	args := []string{
		"--no-playlist",
		"-x",
		"-f", "bestaudio/best",
		"-o", outTmpl,
		opt.URL,
	}

	cmd := exec.CommandContext(ctx, opt.YtDlpPath, args...)
	if err := streamCmd(cmd, log); err != nil {
		return "", fmt.Errorf("yt-dlp download: %w", err)
	}

	// Найдём файл audio.* в tmpDir
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasPrefix(e.Name(), "audio.") {
			return filepath.Join(tmpDir, e.Name()), nil
		}
	}
	return "", errors.New("yt-dlp: audio file not found after download")
}
