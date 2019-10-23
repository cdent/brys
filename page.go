package main

import (
	"io/ioutil"
	"os"
)

type Page struct {
	PageId  string
	Content string
}

// Read a page from disk. `store` is the directory where wiki pages
// are stored.
func (p *Page) read(store string) error {
	content, err := ioutil.ReadFile(p.storeLoc(store))
	if err != nil {
		// FIXME: if the err is not found we should return here. If it is
		// something weird, we should do something about that...
		return err
	}
	p.Content = string(content)
	return nil
}

// Write a page to the store.
// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places.
func (p *Page) save(store string) error {
	f, err := os.Create(p.storeLoc(store))
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

// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places/stores.
func (p *Page) storeLoc(store string) string {
	return store + "/" + p.PageId
}
