package main

import (
	"regexp"
)

// Simple WikiWord, for now, we can't use \b because of unicode
// and golang \b is ascii only (which is fair)
var wikiWord = regexp.MustCompile(`(\A|\s)(\p{Lu}\p{Ll}+\p{Lu}\p{Ll}+)(\p{P}|\s|\z)`)

func wikify(c string) string {
	return wikiWord.ReplaceAllString(c, "$1<a class=\"wikiword\" href=\"/p/$2\">$2</a>$3")
}
