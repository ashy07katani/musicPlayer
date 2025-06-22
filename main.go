package main

import (
	"fmt"
	"music-player/cache"
	handlers "music-player/handlers"
	"music-player/middleware"
	"music-player/repo"
	"net/http"

	//	"os/exec"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db := repo.EstablishConnection()
	redisClient := cache.InitCache()
	h := handlers.NewMusicHandler(db, redisClient)
	r.HandleFunc("/stream/{filename}", h.StreamMusic)
	r.HandleFunc("/upload", h.UploadFile).Methods("POST")
	r.PathPrefix("/chunks/").Handler(http.StripPrefix("/chunks/", http.FileServer(http.Dir("chunks/"))))
	r.HandleFunc("/stream/hls/{filename}", h.StreamHLS).Methods("GET")
	handlerWithCORS := middleware.EnableCORS(r)

	err := http.ListenAndServe(":8080", handlerWithCORS)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("starting server at port :8080")
}
