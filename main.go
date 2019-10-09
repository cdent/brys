package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
)

var box = packr.New("assets", "./assets")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
		s, err := box.FindString("templates/editor.html")
		check(err)
		t, err := template.New("editor").Parse(s)
		check(err)
		data := struct {
			PageId string
		}{
			pageId,
		}
		t.Execute(w, data)
	} else {
		w.Write([]byte("<html><body><h1>something</h1></body></html>"))
	}
}

func setPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	http.Redirect(w, r, fmt.Sprintf("/p/%s", pageId), 303)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/p", func(r chi.Router) {
		r.Get("/{pageId}", getPage)
		r.Post("/{pageId}", setPage)
	})
	r.Get("/", getRoot)

	http.ListenAndServe(":3333", r)
}
