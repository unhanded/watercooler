package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/unhanded/watercooler/internal/msg"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "The port to listen on")
	flag.Parse()

	if port == 0 {
		envPort, err := strconv.ParseInt(os.Getenv("LISTEN_PORT"), 10, 32)
		if err != nil {
			port = 8000
		} else {
			port = int(envPort)
		}
	}

	store := msg.NewMessageStore()

	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		msg := msg.Message{}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Errorf("Error decoding JSON: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := store.Insert(msg)
		if err != nil {
			log.Errorf("Error inserting message: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		b, _ := res.ToJSON(false)
		w.Write(b)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		res, err := store.List()
		if err != nil {
			log.Errorf("Error listing messages: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(res)
		w.Write(b)
	})

	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
