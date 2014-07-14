package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/txgruppi/site/db"
	"github.com/txgruppi/site/links"
	"github.com/txgruppi/site/urlshortener"
)

func main() {
	mongoUrl := os.Getenv("MONGO_URL")
	infoComment := "<!-- " + martini.Env + " -->"

	if mongoUrl == "" {
		panic(errors.New("Can not find environment variable MONGO_URL"))
	}

	mongodb := db.NewMongoDB(mongoUrl)
	defer mongodb.Close()
	err := mongodb.Dial()
	panicIfErr(err)

	database, err := mongodb.Database()
	panicIfErr(err)

	m := martini.Classic()

	linksDAO := links.NewDAO(database.C("links"))
	m.Map(&linksDAO)

	us := urlshortener.New(database.C("urlshortener"))
	m.Map(&us)

	m.Use(render.Renderer(render.Options{
		Layout: "",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "main", map[string]interface{}{"env": template.HTML(infoComment)})
	})

	m.Get("/api/links", func(l *links.LinksDAO) []byte {
		links, err := l.All()
		panicIfErr(err)

		bytes, err := json.Marshal(links)
		panicIfErr(err)

		return bytes
	})

	m.Get("/(?P<id>\\d+)/hashid", func(params martini.Params, us *urlshortener.UrlShortener) string {
		id, err := strconv.Atoi(params["id"])
		if err == nil {
			hashid, err := us.HashIdFor(id)
			if err == nil {
				return hashid
			}
		}
		return "Can't generate hashid"
	})

	m.Get("/:id", func(res http.ResponseWriter, req *http.Request, r render.Render, params martini.Params, us *urlshortener.UrlShortener) {
		url, err := us.UrlFor(params["id"])
		if err == nil && url != "" {
			http.Redirect(res, req, url, 302)
		} else {
			r.HTML(404, "not_found", nil)
		}
	})

	m.Run()
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
