package main

import (
	"fmt"
	"regexp"
)

// Simple WikiWord, for now
var wikiWord = regexp.MustCompile(`\b\p{Lu}\p{Ll}+\p{Lu}\p{Ll}+\b`)

func linkify(w string) string {
	return fmt.Sprintf("<a class=\"wikiword\" href=\"/p/%s\">%s</a>", w, w)
}

func wikify(c string) string {
	return wikiWord.ReplaceAllStringFunc(c, linkify)
}
