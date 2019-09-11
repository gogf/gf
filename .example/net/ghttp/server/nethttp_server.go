package main

import (
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr:         ":8199",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
		time.Sleep(3 * time.Second)
	})
	s.ListenAndServe()
}
