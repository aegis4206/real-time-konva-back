package main

import (
	"log"
	"net/http"
	"real-time-konva-back/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.InitWs)

	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
