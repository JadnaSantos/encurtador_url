package api

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func NewHandler (db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer) 
	r.Use(middleware.RequestID) 
	r.Use(middleware.Logger) 
	r.Use(jsonMiddleware)

	r.Post("/api/shorten", handlePost(db))
	r.Get("/{code}", handleGet(db))

	return r
}

type PostBody struct { 
	URL string `json:"url"`
}

type Respose struct {
	Error string `json:"error,omitempty"`
	Data any `json:"data,omitempty"`
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func sendJSON (w http.ResponseWriter, resp Respose, status int) {
	data, err := json.Marshal(resp) 
	if err != nil {
		slog.Error("failed to marshal json data", "erro", err)
		sendJSON(
			w, 
			Respose{Error: "something went wrong"},
			http.StatusInternalServerError,
		)

		return 
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return 
	}
}

func handlePost (db map[string]string) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w, 
				Respose{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)

			return 
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(
				w,
				Respose{Error: "Invalid url passed"},
				http.StatusBadRequest,
			)
		}

		code := genCode()

		db[code] = body.URL

		sendJSON(
			w,
			Respose{Data: code},
			http.StatusCreated,
		)
}}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genCode () string {
	const n = 8
	byts := make([]byte, n)
	
	for i := range n {
		byts[i] = characters[rand.Intn((len(characters)))]
	}

	return string(byts)
}

func handleGet (db map[string]string) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		data, ok := db[code]

		if !ok {
			http.Error(w, "url not found", http.StatusNotFound)
		}

		http.Redirect(w, r, data, http.StatusPermanentRedirect)
}}

