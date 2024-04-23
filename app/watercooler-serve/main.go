package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sondelll/bulletin/internal/msg"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "The port to listen on")
	flag.Parse()

	if port == 0 {
		envPort, err := strconv.ParseInt(os.Getenv("LISTEN_PORT"), 10, 32)
		if err != nil {
			port = 8880
		} else {
			port = int(envPort)
		}
	}

	store := msg.EchoMessageStore{Store: msg.NewInMemoryMessageStore()}
	ech := echo.New()

	ech.POST("/insert", store.Insert)
	ech.GET("/list", store.List)
	ech.GET("/heartbeat", store.Prune)

	ech.Logger.SetLevel(log.ERROR)

	addr := fmt.Sprintf(":%d", port)

	ech.HideBanner = true
	ech.Start(addr)
}
