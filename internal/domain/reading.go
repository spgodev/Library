package domain

import (
	"time"
)

type ReadingInfo struct {
	User string
	Date time.Time
}

type ReadBookItem struct {
	Book     Book      `json:"book"`
	ReadDate time.Time `json:"read_date"`
}
