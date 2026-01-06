package transcriber

type OpenAI struct {
	APIKey string
	Model  string // gpt-4o-transcribe / whisper-1
}

func (o *OpenAI) Name() string {
	return "openai"
}

func (o *OpenAI) Transcribe(audio string) (*Result, error) {
	// multipart upload → transcription
	// intentionally omitted here (добавим следующим шагом)
	return nil, nil
}
