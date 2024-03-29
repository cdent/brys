package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/packr/v2"
)

var (
	homePage  = flag.String("homepage", "HomePage", "1st Page of wiki")
	pageStore = flag.String("store", "./store", "Path to data storage")
	port      = flag.String("port", "3333", "Port to listen on")
	box       = packr.New("assets", "./assets")
)

// FIXME: this should result in http errors that are dependent on what
// happened, not 500s
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Go to the HomePage when / is requested.
func getRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("/p/%s", *homePage), http.StatusFound)
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
	page.HTML = template.HTML(wikify(page.Content))
	if !page.Modifiedtime.IsZero() {
		w.Header().Set("Last-Modified", page.Modifiedtime.UTC().Format(http.TimeFormat))
	}
	et.Execute(w, page)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	pageId, err := url.QueryUnescape(pageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if pageId == "RecentChanges" {
		sendRecentChanges(w, r)
	} else {
		sendRegularPage(w, r, pageId)
	}
}

func sendRecentChanges(w http.ResponseWriter, r *http.Request) {
	s := &store{base: *pageStore}
	pages, err := listPages(s)
	check(err)
	// this next block seems like there's probably a shortcut that
	// could exist
	b, err := box.FindString("templates/base.html")
	check(err)
	e, err := box.FindString("templates/recents.html")
	check(err)
	bt, err := template.New("page").Parse(b)
	check(err)
	et, err := bt.Parse(e)
	check(err)
	et.Execute(w, pages)
}

func sendRegularPage(w http.ResponseWriter, r *http.Request, pageId string) {
	edit := r.FormValue("edit")

	s := &store{base: *pageStore}
	page := NewPage(pageId, s)
	err := page.read()

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

func setPage(c chan *Page) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageId := chi.URLParam(r, "pageId")
		pageId, err := url.QueryUnescape(pageId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		content := r.PostFormValue("content")
		del := r.PostFormValue("delete")
		if del != "" {
			delPage(w, r)
		} else {
			s := &store{base: *pageStore}
			page := NewPage(pageId, s)
			page.Content = content
			err = page.save()
			check(err)
			// Send a notification of this page.
			page.Modifiedtime = time.Now()
			c <- page
			http.Redirect(w, r, fmt.Sprintf("/p/%s", pageId), http.StatusSeeOther)
		}
	})
}

func delPage(w http.ResponseWriter, r *http.Request) {
	pageId := chi.URLParam(r, "pageId")
	s := &store{base: *pageStore}
	page := NewPage(pageId, s)
	err := page.del()
	check(err)
	http.Redirect(w, r, "/p/HomePage", http.StatusSeeOther)
}

func main() {
	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	// Restrict request bodies to forms (for now)
	r.Use(middleware.AllowContentType("application/x-www-form-urlencoded"))

	hub := newHub()
	go hub.run()

	r.HandleFunc("/ws", serveWs(hub))
	r.Route("/p", func(r chi.Router) {
		r.Get("/{pageId}", getPage)
		r.Post("/{pageId}", setPage(hub.broadcast))
		r.Delete("/{pageId}", delPage)
	})
	r.Get("/", getRoot)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), r))
}
