package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Log struct {
	mutex_ sync.Mutex
	logs   []string
}

var temp_logs = Log{
	logs: []string{},
}

var LogChannel = make(chan string)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	temp_logs.mutex_.Lock()
	w.Write([]byte(strings.Join(temp_logs.logs, "\n")))
	temp_logs.mutex_.Unlock()
}

func listenWatcher(log_channel *chan string) {
	for {
		var get_val string = <-*log_channel
		temp_logs.mutex_.Lock()
		temp_logs.logs = append(temp_logs.logs, get_val)
		temp_logs.mutex_.Unlock()
	}
}

func runServer() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:         "0.0.0.0:5000",
		Handler:      mux,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/getLogs", getLogs)
	fmt.Println("Server starting on port 5000....")
	go listenWatcher(&LogChannel)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
