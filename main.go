package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><h1>Hello World</h1></body></html>"))
}

func readPage(p string) (string, error) {
	// pretend there are no pages, for now
	return "", fmt.Errorf("not found %s", p)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	_, err := readPage(pageId)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("<html><body><h1>%s</h1><form id=\"editor\"></form></body></html>", pageId)))
	} else {
		w.Write([]byte("<html><body><h1>something</h1></body></html>"))
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/p", func(r chi.Router) {
		r.Get("/{pageId}", getPage)
	})
	r.Get("/", getRoot)

	http.ListenAndServe(":3333", r)
}
