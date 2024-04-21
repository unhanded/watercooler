package msg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MessageEchoStore interface {
	Insert(c echo.Context) error
	Find(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Retract(c echo.Context) error
}

type InMemoryMessageStore struct {
	Bulletins []Message
}

func (s *InMemoryMessageStore) Insert(c echo.Context) error {
	b := Message{}
	dec := json.NewDecoder(c.Request().Body)

	if decodeErr := dec.Decode(&b); decodeErr != nil {
		return decodeErr
	}
	// Default fallback for UUID
	if b.Uuid.String() == "" || b.Uuid == uuid.Nil {
		newUuid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.Uuid = newUuid
	}
	// Default fallback for Time of Expiration
	if b.ExpiryUnix == 0 {
		b.ExpiryUnix = time.Now().Add(time.Minute * 30).Unix()
	}
	// Always set the timestamp to the current time, regardless of what the client sends
	b.Timestamp = time.Now().Unix()

	s.Bulletins = append(s.Bulletins, b)

	return nil
}

func (s *InMemoryMessageStore) List(c echo.Context) error {
	byt, err := json.Marshal(s.Bulletins)
	if err != nil {
		return err
	}
	c.Response().Write(byt)
	return nil
}

func NewInMemoryMessageStore() *InMemoryMessageStore {
	return &InMemoryMessageStore{
		Bulletins: []Message{},
	}
}
