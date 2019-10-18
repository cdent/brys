package main

import (
    "os"
	"io/ioutil"
)

type Page struct {
	PageId  string
	Content string
}

// Read a page from disk. `store` is the directory where wiki pages
// are stored.
func readPage(store string, p string) (*Page, error) {
	content, err := ioutil.ReadFile(store + "/" + p)
	if err != nil {
		// FIXME: if the err is not found we should return here. If it is
		// something weird, we should do something about that...
		return nil, err
	}
	return &Page{PageId: p, Content: string(content)}, nil
}

// Write a page to the store.
// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places.
func (p *Page) save(store string) error {
	f, err := os.Create(store + "/" + p.PageId)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(p.Content)
	if err != nil {
		return err
	}
	return nil
}

