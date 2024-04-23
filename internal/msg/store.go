package msg

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MessageStore interface {
	Insert(m Message) (Message, error)
	Find(id string) (Message, error)
	List() ([]Message, error)
	Update(m Message) error
	Retract(id uuid.UUID) error
	Prune() error
}

type InMemoryMessageStore struct {
	Messages []Message
}

func (s *InMemoryMessageStore) Insert(m Message) (Message, error) {
	// Default fallback for UUID
	if m.Uuid.String() == "" || m.Uuid == uuid.Nil {
		newUuid, err := uuid.NewV7()
		if err != nil {
			return Message{}, err
		}
		m.Uuid = newUuid
	}
	// Default fallback for Time of Expiration
	if m.ExpiryUnix == 0 {
		m.ExpiryUnix = time.Now().Add(time.Minute * 30).Unix()
	}
	// Always set the timestamp to the current time, regardless of what the client sends
	m.Timestamp = time.Now().Unix()

	s.Messages = append(s.Messages, m)

	return m, nil
}

func (s *InMemoryMessageStore) Find(id string) (Message, error) {
	return Message{}, fmt.Errorf("not implemented")
}

func (s *InMemoryMessageStore) List() ([]Message, error) {
	return s.Messages, nil
}

func (s *InMemoryMessageStore) Update(m Message) error {
	return fmt.Errorf("not implemented")
}

func (s *InMemoryMessageStore) Prune() error {
	for i, b := range s.Messages {
		if b.ExpiryUnix < time.Now().Unix() {
			s.Messages = append(s.Messages[:i], s.Messages[i+1:]...)
		}
	}
	return nil
}

func (s *InMemoryMessageStore) Retract(id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func NewInMemoryMessageStore() *InMemoryMessageStore {
	return &InMemoryMessageStore{
		Messages: []Message{},
	}
}
