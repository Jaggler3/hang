package chat

import "time"

type Message struct {
	PlayerName string
	Text       string
	Time       time.Time
}

type Buffer struct {
	Messages []Message
	Offset   int
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Add(name, text string) {
	b.Messages = append(b.Messages, Message{PlayerName: name, Text: text, Time: time.Now()})
}

func (b *Buffer) Since(offset int) ([]Message, int) {
	if offset >= len(b.Messages) {
		return nil, len(b.Messages)
	}
	return b.Messages[offset:], len(b.Messages)
}
