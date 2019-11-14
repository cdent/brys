package main

import (
	"fmt"
	"testing"
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

func TestFullStringWikify(t *testing.T) {
	lines := [][]string{
		{"We have OneWord and TwoWord in this line.", "We have <a class=\"wikiword\" href=\"/p/OneWord\">OneWord</a> and <a class=\"wikiword\" href=\"/p/TwoWord\">TwoWord</a> in this line."},
		{"Now OneWord TwoWord next to one another.", "Now <a class=\"wikiword\" href=\"/p/OneWord\">OneWord</a> <a class=\"wikiword\" href=\"/p/TwoWord\">TwoWord</a> next to one another."},
		{"OneWord TwoWord at the start.", "<a class=\"wikiword\" href=\"/p/OneWord\">OneWord</a> <a class=\"wikiword\" href=\"/p/TwoWord\">TwoWord</a> at the start."},
		{"At the end is OneWord TwoWord", "At the end is <a class=\"wikiword\" href=\"/p/OneWord\">OneWord</a> <a class=\"wikiword\" href=\"/p/TwoWord\">TwoWord</a>"},
		{"At the end is OneWord TwoWord.", "At the end is <a class=\"wikiword\" href=\"/p/OneWord\">OneWord</a> <a class=\"wikiword\" href=\"/p/TwoWord\">TwoWord</a>."},
	}

	for _, line := range lines {
		wiki := wikify(line[0])
		if wiki != line[1] {
			t.Errorf("\"%s\" wikified to \"%s\", expected \"%s\"", line[0], wiki, line[1])
		}
	}
}
