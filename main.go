package main

import (
	"fmt"
	handlers "music-player/Handlers"
	"music-player/middleware"
	"net/http"

	//	"os/exec"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stream/{filename}", handlers.StreamMusic)
	r.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	r.PathPrefix("/chunks/").Handler(http.StripPrefix("/chunks/", http.FileServer(http.Dir("chunks/"))))
	r.HandleFunc("/stream/hls/{filename}", handlers.StreamHLS).Methods("GET")
	handlerWithCORS := middleware.EnableCORS(r)
	err := http.ListenAndServe(":8080", handlerWithCORS)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("starting server at port :8080")
}
