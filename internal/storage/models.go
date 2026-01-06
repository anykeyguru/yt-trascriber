package storage

import "time"

type Transcript struct {
	ID        int64
	Title     string
	URL       string
	Backend   string
	Text      string
	CreatedAt time.Time
}
