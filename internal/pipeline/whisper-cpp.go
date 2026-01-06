package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func whisperCPP(ctx context.Context, opt Options, wavPath, title string, log LogFn) (string, error) {
	if opt.WhisperBin == "" {
		opt.WhisperBin = "whisper-cli"
	}
	if opt.WhisperModel == "" {
		return "", errors.New("whisper.cpp: model path is empty")
	}

	base := CreateFileNameForSave(title)
	if base == "" {
		base = "transcript"
	}
	outPrefix := filepath.Join(opt.OutDir, base)

	// whisper-cli обычно создаёт outPrefix + ".txt" при -otxt
	args := []string{
		"-m", opt.WhisperModel,
		"-f", wavPath,
		"-of", outPrefix,
		"-otxt",
	}
	if opt.Language != "" && opt.Language != "auto" {
		args = append(args, "-l", opt.Language)
	}

	cmd := exec.CommandContext(ctx, opt.WhisperBin, args...)
	if err := streamCmd(cmd, log); err != nil {
		return "", fmt.Errorf("whisper.cpp: %w", err)
	}

	txtPath := outPrefix + ".txt"
	b, err := os.ReadFile(txtPath)
	if err != nil {
		// Иногда whisper.cpp пишет в рабочий каталог, но мы даём -of полный путь — должно быть ок.
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("whisper.cpp did not produce %s", txtPath)
		}
		return "", err
	}

	log("saved: " + txtPath)
	return string(b), nil
}
