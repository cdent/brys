package main

import (
	"io/ioutil"
	"os"
	"sync"
)

type Page struct {
	PageId  string
	Content string
	Store   *store
	// Mutex around read operations so that we don't read or write a
	// page while it is already locked for writing. Multiple readers
	// are fine.
	mutex sync.RWMutex
}

// FIXME: This will become more intereting with time,
// but for now we just want something to occupy the slot.
type store struct {
	base string
}

// Read a page from disk.
func (p *Page) read() error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	content, err := ioutil.ReadFile(p.storeLoc())
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
func (p *Page) save() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	f, err := os.Create(p.storeLoc())
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

func (p *Page) del() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	err := os.Remove(p.storeLoc())
	return err
}

// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places/stores.
func (p *Page) storeLoc() string {
	return p.Store.base + "/" + p.PageId
}
