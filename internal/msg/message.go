package msg

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	Uuid        uuid.UUID `json:"uuidV7"`
	Author      string    `json:"author"`
	Topic       []string  `json:"topic"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	LifetimeSec int       `json:"lifetimeSec"`
	Timestamp   int64     `json:"timestamp"`
}

func (b *Message) ToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(b, "", "  ")
	} else {
		return json.Marshal(b)
	}
}

func (b *Message) FromJSON(data []byte) error {
	return json.Unmarshal(data, b)
}
