package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><h1>Hello World</h1></body></html>"))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", getRoot)

	http.ListenAndServe(":3333", r)
}
