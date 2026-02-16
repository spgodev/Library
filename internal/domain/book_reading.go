package domain

import "time"

type BookReading struct {
	ID       int64
	BookID   int64
	UserID   int64
	ReadDate time.Time
}
