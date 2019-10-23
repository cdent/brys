package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
)

var (
	homePage = "HomePage"
	store    = "./store"
	box      = packr.New("assets", "./assets")
)

// FIXME: this should result in http errors, not fatal exits
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Go to the HomePage when / is requested.
func getRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("/p/%s", homePage), http.StatusFound)
}

func sendEditor(w http.ResponseWriter, page *Page) {
	b, err := box.FindString("templates/base.html")
	check(err)
	e, err := box.FindString("templates/editor.html")
	check(err)
	bt, err := template.New("editor").Parse(b)
	check(err)
	et, err := bt.Parse(e)
	check(err)
	et.Execute(w, page)
}

func sendWikiPage(w http.ResponseWriter, page *Page) {
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
		page.PageId,
		// FIXME: This allows anything HTML in the wikification to pass through
		template.HTML(wikify(page.Content)),
	}
	et.Execute(w, data)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	pageId, err := url.QueryUnescape(pageId)
	check(err)
	edit := r.FormValue("edit")

	page := &Page{PageId: pageId}
	err = page.read(store)

	// If we tried to get the page and it is not there, create
	// a new one in the editor.
	if errors.Is(err, os.ErrNotExist) {
		edit = "edit"
	} else {
		check(err)
	}

	if edit == "edit" {
		sendEditor(w, page)
	} else {
		sendWikiPage(w, page)
	}
}

func setPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	pageId, err := url.QueryUnescape(pageId)
	check(err)
	content := r.PostFormValue("content")
	page := &Page{PageId: pageId, Content: content}
	err = page.save(store)
	check(err)
	http.Redirect(w, r, fmt.Sprintf("/p/%s", pageId), http.StatusSeeOther)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/p", func(r chi.Router) {
		r.Get("/{pageId}", getPage)
		r.Post("/{pageId}", setPage)
	})
	r.Get("/", getRoot)

	http.ListenAndServe(":3333", r)
}
