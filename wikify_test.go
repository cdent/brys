
package main

import (
    "testing"
    "fmt"
)

func TestWikify(t *testing.T) {
    rus := []string{
        "CowBar",
        "АppleЯow",
    }

    for _, v := range rus {
        wiki := wikify(fmt.Sprintf("Hi %s Hello", v))
        if wiki != fmt.Sprintf("Hi <a class=\"wikiword\" href=\"/p/%s\">%s</a> Hello", v, v) {
            t.Errorf("%s resulted in %s", v, wiki)
        }
    }
}
