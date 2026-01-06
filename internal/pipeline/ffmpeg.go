package pipeline

import (
	"context"
	"fmt"
	"os/exec"
)

func ffmpegToWav16kMono(ctx context.Context, opt Options, in, out string, log LogFn) error {
	args := []string{
		"-y",
		"-i", in,
		"-ac", "1",
		"-ar", "16000",
		"-vn",
		out,
	}
	cmd := exec.CommandContext(ctx, opt.FfmpegPath, args...)
	if err := streamCmd(cmd, log); err != nil {
		return fmt.Errorf("ffmpeg convert: %w", err)
	}
	return nil
}
