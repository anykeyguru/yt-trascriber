package transcriber

import "os/exec"

type WhisperCPP struct {
	Bin   string
	Model string
	Lang  string
}

func (w *WhisperCPP) Name() string {
	return "whisper.cpp"
}

func (w *WhisperCPP) Transcribe(audio string) (*Result, error) {
	cmd := exec.Command(
		w.Bin,
		"-m", w.Model,
		"-f", audio,
		"-l", w.Lang,
		"-otxt",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return &Result{
		Text: string(out),
	}, nil
}
