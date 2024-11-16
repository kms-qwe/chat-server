package model

import "time"

// Message модель для работы с postgres
type Message struct {
	From     string    `db:"user_name"`
	Text     string    `db:"message_text"`
	ChatID   int64     `db:"chat_id"`
	SendTime time.Time `db:"message_time_send"`
}
