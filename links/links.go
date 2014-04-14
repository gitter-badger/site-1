package links

import (
	"encoding/json"
	mdc "github.com/txgruppi/modtimechecker"
	"html/template"
	"os"
)

type Links struct {
	checker mdc.Checker
	links   []Link
}

type Link struct {
	Title   string
	Url     string
	SafeUrl template.URL
	Image   string
}

func (l *Links) Links() []Link {
	l.checker.Check()
	return l.links
}

func New(path string) *Links {
	links := Links{}

	links.checker = mdc.New(path, func(path string) {
		file, err := os.Open(path)
		if err != nil {
			return
		}
		defer file.Close()

		dec := json.NewDecoder(file)
		dec.Decode(&links.links)
	})

	return &links
}
