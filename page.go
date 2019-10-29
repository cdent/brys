package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var pageLockMutex sync.Mutex

// FIXME: This will grow to the total number of pages.
var pageLocks = make(map[string]*sync.RWMutex)

type Page struct {
	PageId       string
	Content      string
	HTML         template.HTML
	Modifiedtime time.Time
	Store        *store
}

// FIXME: This will become more intereting with time,
// but for now we just want something to occupy the slot.
type store struct {
	base string
}

func NewPage(title string, s *store) *Page {
	page := &Page{PageId: title, Content: "", Store: s}
	return page
}

func listPages() []Page {
	nilStore := &store{base: ""}
	page := Page{PageId: "HomePage", Content: "", Store: nilStore}
	return []Page{page}
}

// Read a page from disk.
func (p *Page) read() error {
	locker := p.readLock()
	defer locker.RUnlock()
	content, err := ioutil.ReadFile(p.storeLoc())
	if err != nil {
		// FIXME: if the err is not found we should return here. If it is
		// something weird, we should do something about that...
		return err
	}
	p.Content = string(content)
	file, err := os.Stat(p.storeLoc())
	// Ignore err if there is one.
	if err == nil {
		p.Modifiedtime = file.ModTime()
	}
	return nil
}

// Write a page to the store.
// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places.
func (p *Page) save() error {
	locker := p.writeLock()
	defer locker.Unlock()
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
	locker := p.writeLock()
	defer locker.Unlock()
	err := os.Remove(p.storeLoc())
	return err
}

func (p *Page) writeLock() *sync.RWMutex {
	pageLockMutex.Lock()
	locker, ok := pageLocks[p.PageId]
	if !ok {
		locker = &sync.RWMutex{}
		pageLocks[p.PageId] = locker
	}
	locker.Lock()
	pageLockMutex.Unlock()
	return locker
}

func (p *Page) readLock() *sync.RWMutex {
	pageLockMutex.Lock()
	locker, ok := pageLocks[p.PageId]
	if !ok {
		locker = &sync.RWMutex{}
		pageLocks[p.PageId] = locker
	}
	locker.RLock()
	pageLockMutex.Unlock()
	return locker
}

// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places/stores.
func (p *Page) storeLoc() string {
	return p.Store.base + "/" + p.PageId
}
