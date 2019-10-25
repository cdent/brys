// Test various page routines not covered by gabbi tests.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTwoStores(t *testing.T) {
	t1, err := ioutil.TempDir("", "s1")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(t1)
	t2, err := ioutil.TempDir("", "s2")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(t2)
	store1 := &store{base: t1}
	store2 := &store{base: t2}
	page1 := &Page{PageId: "1", Content: "1", Store: store1}
	page2 := &Page{PageId: "2", Content: "2", Store: store2}
	page1.save()
	page2.save()

	if page1.Store.base != t1 {
		t.Errorf("page1 has the wrong Store base: %s", page1.Store.base)
	}
	if page2.Store.base != t2 {
		t.Errorf("page2 has the wrong Store base: %s", page2.Store.base)
	}

	content1, err := ioutil.ReadFile(t1 + "/1")
	if err != nil {
		log.Fatal(err)
	}
	content2, err := ioutil.ReadFile(t2 + "/2")
	if err != nil {
		log.Fatal(err)
	}

	if string(content1) != "1" {
		t.Errorf("content1 should contain 1 but has %s", content1)
	}
	if string(content2) != "2" {
		t.Errorf("content2 should contain 2 but has %s", content2)
	}
}

func TestPageMutex(t *testing.T) {
	s := &store{base: "."}
	p1 := NewPage("1", s)
	p1.Content = "1"
	p2 := NewPage("1", s)
	p2.Content = "2"
	p3 := NewPage("1", s)
	p3.Content = "3"

	c := make(chan string, 1)
	go func() {
		pageLockMutex.Lock()
		p1Lock, ok := pageLocks[p1.PageId]
		if !ok {
			p1Lock = &sync.RWMutex{}
			pageLocks[p1.PageId] = p1Lock
		}
		p1Lock.Lock()
		pageLockMutex.Unlock()
		c <- "p1 locked"
	}()

	go func() {
		// make sure p1 goes first
		time.Sleep(5 * time.Millisecond)
		pageLockMutex.Lock()
		p2Lock, ok := pageLocks[p2.PageId]
		if !ok {
			p2Lock = &sync.RWMutex{}
			pageLocks[p2.PageId] = p2Lock
		}
		p2Lock.Lock()
		pageLockMutex.Unlock()
		c <- "p2 locked"
	}()

	p1Locked := false
	p2Locked := false

outer:
	for {
		select {
		case res := <-c:
			t.Log("res", res)
			if res == "p1 locked" {
				p1Locked = true
			}
			if res == "p2 locked" {
				p2Locked = true
			}

		case <-time.After(25 * time.Millisecond):
			break outer
		}
	}

	if p2Locked {
		t.Log("p2 locked when it should not have")
		t.Fail()
	}
	if !p1Locked {
		t.Log("p1 never locked")
		t.Fail()
	}
}
