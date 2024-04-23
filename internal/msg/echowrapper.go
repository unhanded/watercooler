package msg

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type EchoStoreWrapper interface {
	Insert(c echo.Context) error
	Find(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Retract(c echo.Context) error
	Prune(c echo.Context) error
}

type EchoMessageStore struct {
	Store MessageStore
}

func (e *EchoMessageStore) Insert(c echo.Context) error {
	var m Message = Message{}
	dec := json.NewDecoder(c.Request().Body)
	desErr := dec.Decode(&m)
	if desErr != nil {
		return desErr
	}
	inserted, err := e.Store.Insert(m)
	if err != nil {
		return c.String(500, "An error occured while inserting the message")
	}

	return c.JSON(200, inserted)
}

func (e *EchoMessageStore) List(c echo.Context) error {
	all, err := e.Store.List()
	if err != nil {
		return c.String(500, "An error occured while listing the messages")
	}
	return c.JSON(200, all)
}

func (e *EchoMessageStore) Prune(c echo.Context) error {
	err := e.Store.Prune()
	if err != nil {
		return c.String(500, "Internal server error")
	}
	return c.String(200, "ACK")
}
