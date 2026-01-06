package transcriber

type Result struct {
	Text     string
	Language string
	Segments []Segment
}

type Segment struct {
	Start float64
	End   float64
	Text  string
}

type Transcriber interface {
	Name() string
	Transcribe(audioPath string) (*Result, error)
}
