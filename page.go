package main

import (
	"io/ioutil"
	"os"
	"sync"
)

var pageLockMutex sync.Mutex

// FIXME: This will grow to the total number of pages.
var pageLocks = make(map[string]*sync.RWMutex)

type Page struct {
	PageId  string
	Content string
	Store   *store
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

// Read a page from disk.
func (p *Page) read() error {
	pageLockMutex.Lock()
	locker, ok := pageLocks[p.PageId]
	if !ok {
		locker = &sync.RWMutex{}
		pageLocks[p.PageId] = locker
	}
	locker.RLock()
	pageLockMutex.Unlock()
	defer locker.RUnlock()
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
	pageLockMutex.Lock()
	locker, ok := pageLocks[p.PageId]
	if !ok {
		locker = &sync.RWMutex{}
		pageLocks[p.PageId] = locker
	}
	locker.Lock()
	pageLockMutex.Unlock()
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
	pageLockMutex.Lock()
	locker, ok := pageLocks[p.PageId]
	if !ok {
		locker = &sync.RWMutex{}
		pageLocks[p.PageId] = locker
	}
	locker.Lock()
	pageLockMutex.Unlock()
	defer locker.Unlock()
	err := os.Remove(p.storeLoc())
	return err
}

// FIXME: consider having store as an attribute on the Page, so pages could
// be in different places/stores.
func (p *Page) storeLoc() string {
	return p.Store.base + "/" + p.PageId
}
