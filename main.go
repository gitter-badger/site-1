package main

import (
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/txgruppi/site/links"
	"github.com/txgruppi/site/urlshortener"
	"html/template"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	targetDate, _ := time.Parse("2006-01-02 15:04:05", "2014-05-03 20:00:00")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	m := martini.Classic()

	l := links.New(path.Join(cwd, "data", "links.json"))

	us := urlshortener.New("urlshortener")
	if us == nil {
		panic(errors.New("Can't initialize UrlShortener"))
	}

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		now := time.Now().Add(-(3 * time.Hour))
		daysLeft := int(targetDate.Sub(now).Hours() / 24)
		if daysLeft < 0 {
			daysLeft = 0
		}

		links := l.Links()

		for key := range links {
			if links[key].Title == "Skype" && links[key].Url[0:8] == "skype://" {
				links[key].SafeUrl = template.URL(links[key].Url)
			}
		}

		r.HTML(200, "index", map[string]interface{}{"links": links, "daysLeft": daysLeft})
	})

	m.Get("/:id", func(res http.ResponseWriter, req *http.Request, r render.Render, params martini.Params) {
		url, err := us.UrlFor(params["id"])
		if err == nil && url != "" {
			http.Redirect(res, req, url, 302)
		} else {
			r.HTML(404, "not_found", nil)
		}
	})

	m.NotFound(func(r render.Render) {
		r.HTML(404, "not_found", nil)
	})

	m.Run()
}
