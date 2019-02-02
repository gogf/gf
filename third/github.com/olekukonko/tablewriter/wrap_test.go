// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"strings"
	"testing"

	"github.com/gogf/gf/third/github.com/mattn/go-runewidth"
)

var text = "The quick brown fox jumps over the lazy dog."

func TestWrap(t *testing.T) {
	exp := []string{
		"The", "quick", "brown", "fox",
		"jumps", "over", "the", "lazy", "dog."}

	got, _ := WrapString(text, 6)
	checkEqual(t, len(got), len(exp))
}

func TestWrapOneLine(t *testing.T) {
	exp := "The quick brown fox jumps over the lazy dog."
	words, _ := WrapString(text, 500)
	checkEqual(t, strings.Join(words, string(sp)), exp)

}

func TestUnicode(t *testing.T) {
	input := "Česká řeřicha"
	var wordsUnicode []string
	if runewidth.IsEastAsian() {
		wordsUnicode, _ = WrapString(input, 14)
	} else {
		wordsUnicode, _ = WrapString(input, 13)
	}
	// input contains 13 (or 14 for CJK) runes, so it fits on one line.
	checkEqual(t, len(wordsUnicode), 1)
}

func TestDisplayWidth(t *testing.T) {
	input := "Česká řeřicha"
	want := 13
	if runewidth.IsEastAsian() {
		want = 14
	}
	if n := DisplayWidth(input); n != want {
		t.Errorf("Wants: %d Got: %d", want, n)
	}
	input = "\033[43;30m" + input + "\033[00m"
	checkEqual(t, DisplayWidth(input), want)
}
