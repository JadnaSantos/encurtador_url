package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/JadnaSantos/Encurtador-de-Url/api"
)

func main () {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err) 
		return 
	}

	slog.Info("all systems offline")
}


func run () error {
	db := make(map[string]string)

	handle := api.NewHandler(db)

	s := http.Server{
		Addr:                         ":8080",
		Handler:                      handle,	
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}