package golem

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"path"
	"strings"

	"github.com/bountylabs/golem/dicts"
)

// Lemmatizer is the key to lemmatizing a word in a language
type Lemmatizer struct {
	m map[string][]string
}

const folder = "data"

// New produces a new Lemmatizer
func New(locale string) (*Lemmatizer, error) {
	var fname string

	switch locale {
	case "sv", "swedish":
		fname = "sv.gz"
	case "en", "english":
		fname = "en.gz"
	default:
		return nil, fmt.Errorf(`Language "%s" is not implemented`, locale)
	}
	f, err := dicts.Asset(path.Join(folder, fname))
	if err != nil {
		return nil, fmt.Errorf(`Could not open resource file for "%s"`, locale)
	}
	r, err := gzip.NewReader(bytes.NewBuffer(f))
	if err != nil {
		return nil, fmt.Errorf(`Could not open resource file for "%s"`, locale)
	}

	l := Lemmatizer{m: make(map[string][]string)}
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	for err == nil {
		parts := strings.Split(strings.TrimSpace(line), "\t")
		if len(parts) == 2 {
			base := strings.ToLower(parts[0])
			form := strings.ToLower(parts[1])
			l.m[form] = append(l.m[form], base)
		} else {
			fmt.Println(line, "is odd")
		}
		line, err = br.ReadString('\n')
	}
	return &l, nil
}

// Lemma gets one of the base forms of a word
func (l *Lemmatizer) Lemma(word string) (string, error) {
	if out, ok := l.m[strings.ToLower(word)]; ok {
		return out[0], nil
	}
	return "", fmt.Errorf("Word not found in dictionary")
}

// Lemmas gets all the base forms of a word
func (l *Lemmatizer) Lemmas(word string) ([]string, error) {
	if out, ok := l.m[strings.ToLower(word)]; ok {
		return out, nil
	}
	return nil, fmt.Errorf("Word not found in dictionary")
}
