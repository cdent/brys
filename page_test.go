// Test various page routines not covered by gabbi tests.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
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
