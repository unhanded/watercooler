package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sondelll/bulletin/internal/msg"
)

func main() {
	listenPort := os.Getenv("LISTEN_PORT")

	store := msg.NewInMemoryMessageStore()
	ech := echo.New()

	ech.POST("/insert", store.Insert)
	ech.GET("/list", store.List)

	ech.Logger.SetLevel(log.ERROR)

	addr := fmt.Sprintf(":%s", listenPort)
	ech.Start(addr)
}
