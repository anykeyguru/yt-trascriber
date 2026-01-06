Отлично. Это очень разумный и расширяемый дизайн. Ниже — **чёткая архитектура**, **стек библиотек**, **структура проекта** и **рабочий TUI-скелет на Bubble Tea**, уже с:

* выбором backend: **local whisper.cpp / OpenAI API**
* пайплайном: YouTube → audio → transcript
* **SQLite-хранилищем результатов**
* **экспортом результата в файл по указанному пути**
* без перехвата clipboard (как ты и запросил)

Я deliberately держу код **production-ориентированным**, а не «демо ради демо».

---

# 1. Архитектура (high-level)

```
┌───────────┐
│ Bubble TUI│
│ (bubbletea)
└─────┬─────┘
      │ commands / events
┌─────▼────────┐
│   Controller │
│ (job manager)│
└─────┬────────┘
      │
┌─────▼─────────────────────────────────────┐
│               Pipeline                     │
│  yt-dlp → ffmpeg → transcriber (interface) │
└─────┬───────────────────────────┬─────────┘
      │                           │
┌─────▼────────────┐     ┌────────▼──────────┐
│ Whisper.cpp impl │     │ OpenAI API impl   │
└──────────────────┘     └───────────────────┘
      │
┌─────▼──────────┐
│ SQLite storage │
│ transcripts    │
└─────┬──────────┘
      │
┌─────▼─────────────┐
│ Export (TXT/SRT…) │
└───────────────────┘
```

Ключевая идея: **backend транскрибации — интерфейс**, TUI ничего о нём не знает.

---

# 2. Технологический стек (MacBook-friendly)

### Go libs

* `github.com/charmbracelet/bubbletea`
* `github.com/charmbracelet/bubbles`
* `github.com/charmbracelet/lipgloss`
* `modernc.org/sqlite` (чистый Go, без CGO)
* `github.com/spf13/cobra` (CLI entrypoint, опционально)

### External tools

* `yt-dlp`
* `ffmpeg`
* `whisper.cpp` **или** OpenAI API

---

# 3. Структура проекта (рекомендую именно так)

```
ytranscribe/
├── cmd/
│   └── ytranscribe/
│       └── main.go
├── internal/
│   ├── tui/
│   │   ├── model.go
│   │   ├── view.go
│   │   ├── update.go
│   │   └── keys.go
│   ├── pipeline/
│   │   ├── pipeline.go
│   │   ├── downloader.go
│   │   ├── audio.go
│   │   └── job.go
│   ├── transcriber/
│   │   ├── transcriber.go   // interface
│   │   ├── whispercpp.go
│   │   └── openai.go
│   ├── storage/
│   │   ├── sqlite.go
│   │   └── models.go
│   └── export/
│       └── export.go
├── migrations/
│   └── 001_init.sql
├── Makefile
└── go.mod
```

---

# 4. Интерфейс транскрибации (ключевая точка расширения)

```go
// internal/transcriber/transcriber.go
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
```

---

## Whisper.cpp реализация (локально)

```go
// internal/transcriber/whispercpp.go
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
```

---

## OpenAI API реализация (fallback)

```go
// internal/transcriber/openai.go
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
```

---

# 5. SQLite хранилище

### Миграция

```sql
-- migrations/001_init.sql
CREATE TABLE transcripts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT,
  source_url TEXT,
  backend TEXT,
  language TEXT,
  text TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Репозиторий

```go
// internal/storage/sqlite.go
package storage

import "database/sql"

type Store struct {
	db *sql.DB
}

func (s *Store) Save(t Transcript) error {
	_, err := s.db.Exec(`
		INSERT INTO transcripts (title, source_url, backend, language, text)
		VALUES (?, ?, ?, ?, ?)`,
		t.Title, t.URL, t.Backend, t.Language, t.Text,
	)
	return err
}
```

---

# 6. Экспорт в файл (любой путь)

```go
// internal/export/export.go
package export

import "os"

func ToFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
```

TUI просто спрашивает путь → вызывает `export.ToFile(...)`.

---

# 7. Bubble Tea: модель (минимальный, но правильный старт)

```go
// internal/tui/model.go
package tui

import "github.com/charmbracelet/bubbletea"

type JobStatus int

const (
	Idle JobStatus = iota
	Downloading
	Converting
	Transcribing
	Done
	Error
)

type Job struct {
	Title  string
	URL    string
	Status JobStatus
	Logs   []string
}

type Model struct {
	Jobs        []Job
	Selected    int
	Backend     string // whisper.cpp | openai
	Message     string
}

func InitialModel() Model {
	return Model{
		Backend: "whisper.cpp",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
```

---

# 8. Управление (keymap, k9s-style)

Пример:

* `a` — add URL
* `b` — switch backend
* `enter` — run job
* `s` — save transcript to file
* `q` — quit

Bubble Tea позволяет это сделать очень аккуратно.

---

# 9. Почему это **хороший продукт**, а не скрипт

* backend-agnostic (локально / cloud)
* воспроизводимость (SQLite)
* контроль данных (локальный Whisper)
* масштабируемо (batch, history, re-export)
* TUI → идеально для developer workflow

По факту — это **k9s для транскрибации**.

---

# 10. Следующий логичный шаг

Я предлагаю дальше идти **пошагово**, без перегруза:

**Шаг 1**
→ Я дам тебе **полностью рабочий Bubble Tea TUI** (jobs list + logs + backend switch)

**Шаг 2**
→ Добавим SQLite + history view

**Шаг 3**
→ Подключим OpenAI backend как fallback

**Шаг 4**
→ GoReleaser + single binary for macOS

Если согласен — скажи:

> «Давай шаг 1: TUI каркас»

И я сразу дам готовый `internal/tui` с рабочим `main.go`.
