package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
)

var homePage = "HomePage"
var store = "./store"
var box = packr.New("assets", "./assets")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("/p/%s", homePage), http.StatusFound)
}

func readPage(p string) (string, error) {
	content, err := ioutil.ReadFile(store + "/" + p)
	if err != nil {
		return "", fmt.Errorf("not found %s", p)
	}
	return string(content), nil
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	edit := r.FormValue("edit")
	content, err := readPage(pageId)
	if err != nil || edit == "edit" {
		b, err := box.FindString("templates/base.html")
		check(err)
		e, err := box.FindString("templates/editor.html")
		check(err)
		bt, err := template.New("editor").Parse(b)
		check(err)
		et, err := bt.Parse(e)
		check(err)
		data := struct {
			PageId  string
			Content string
		}{
			pageId,
			content,
		}
		et.Execute(w, data)
	} else {
		b, err := box.FindString("templates/base.html")
		check(err)
		e, err := box.FindString("templates/page.html")
		check(err)
		bt, err := template.New("page").Parse(b)
		check(err)
		et, err := bt.Parse(e)
		check(err)
		data := struct {
			PageId  string
			Content template.HTML
		}{
			pageId,
			// FIXME: This allows anything HTML in the wikification to pass through
			template.HTML(wikify(content)),
		}
		et.Execute(w, data)
	}
}

func setPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	content := r.PostFormValue("content")
	f, err := os.Create(store + "/" + pageId)
	check(err)
	defer f.Close()
	_, err = f.WriteString(content)
	check(err)
	f.Sync()
	http.Redirect(w, r, fmt.Sprintf("/p/%s", pageId), http.StatusSeeOther)
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
