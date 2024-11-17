package model

import "time"

// Message holds information about a message.
type Message struct {
	From     string
	Text     string
	ChatID   int64
	SendTime time.Time
}
