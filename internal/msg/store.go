package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	defer s.Prune()
	// Default fallback for UUID
	if m.Uuid.String() == "" || m.Uuid == uuid.Nil {
		newUuid, err := uuid.NewV7()
		if err != nil {
			return Message{}, err
		}
		m.Uuid = newUuid
	}
	// Default fallback for lifetime if none is provided
	if m.LifetimeSec == 0 {
		m.LifetimeSec = 60 * 30 // 30 minutes
	}
	// Always set the timestamp to the current time, regardless of what the client sends
	m.Timestamp = time.Now().Unix()

	s.Messages = append(s.Messages, m)

	return m, nil
}

func (s *InMemoryMessageStore) Find(id string) (Message, error) {
	s.Prune()
	return Message{}, fmt.Errorf("not implemented")
}

func (s *InMemoryMessageStore) List() ([]Message, error) {
	return s.Messages, nil
}

func (s *InMemoryMessageStore) Update(m Message) error {
	s.Prune()
	return fmt.Errorf("not implemented")
}

func (s *InMemoryMessageStore) Prune() error {
	for i, b := range s.Messages {
		if b.Timestamp+int64(b.LifetimeSec) < time.Now().Unix() {
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

type ValkeyMessageStore struct {
	Conn *redis.Conn
}

func (s *ValkeyMessageStore) Insert(m Message) (Message, error) {
	// Default fallback for UUID
	if m.Uuid.String() == "" || m.Uuid == uuid.Nil {
		newUuid, err := uuid.NewV7()
		if err != nil {
			return Message{}, err
		}
		m.Uuid = newUuid
	}
	// Default fallback for lifetime if none is provided
	if m.LifetimeSec == 0 {
		m.LifetimeSec = int(60 * 30) // 30 minutes
	}
	// Always set the timestamp to the current time, regardless of what the client sends
	m.Timestamp = time.Now().Unix()

	// Set the message in the Redis store
	valkeyExpiration := time.Duration(m.LifetimeSec) * time.Second
	val, jErr := m.ToJSON(false)
	if jErr != nil {
		return Message{}, jErr
	}

	_, err := s.Conn.Set(context.TODO(), m.Uuid.String(), []byte(val), valkeyExpiration).Result()
	if err != nil {
		return Message{}, err
	}

	return m, nil
}

func (s *ValkeyMessageStore) Find(id string) (Message, error) {
	val, err := s.Conn.Get(context.TODO(), id).Result()
	if err != nil {
		return Message{}, err
	}
	var m Message
	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		return Message{}, err
	}
	return m, nil
}

func (s *ValkeyMessageStore) List() ([]Message, error) {
	vals, err := s.Conn.Keys(context.TODO(), "*").Result()
	if err != nil {
		return []Message{}, err
	}
	var messages []Message
	for _, val := range vals {
		m, err := s.Find(val)
		if err != nil {
			return []Message{}, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *ValkeyMessageStore) Prune() error {
	return nil // Valkey handles this for us
}

func (s *ValkeyMessageStore) Update(m Message) error {
	return fmt.Errorf("not implemented")
}

func (s *ValkeyMessageStore) Retract(id uuid.UUID) error {
	_, err := s.Conn.Del(context.TODO(), id.String()).Result()
	return err
}

func NewValkeyMessageStore() MessageStore {
	store := ValkeyMessageStore{}

	client := redis.NewClient(&redis.Options{
		Addr: redisEnv(),
	})

	store.Conn = client.Conn()
	if store.Conn == nil {
		fmt.Println("Failed to connect to Redis, falling back to in-memory store")
		return NewInMemoryMessageStore()
	}
	return &store
}

func redisEnv() string {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	return fmt.Sprintf("%s:%s", redisHost, redisPort)
}
