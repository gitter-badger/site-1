package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"log"

	"labix.org/v2/mgo"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/txgruppi/site/db"
	"github.com/txgruppi/site/links"
	"github.com/txgruppi/site/urlshortener"
	metrics "github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
)

func main() {
	mongoUrl := os.Getenv("MONGO_URL")

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

	registerNewRelic(m)
	registerLinksDao(m, database)
	registerUrlShortener(m, database)
	registerRenderer(m)

	m.Get("/", getRoot)
	m.Get("/api/links", getLinks)
	m.Get("/(?P<id>\\d+)/hashid", showHashidForId)
	m.Get("/:id", redirectToUrl)

	m.Run()
}

func registerNewRelic(m *martini.ClassicMartini) {
	license := os.Getenv("NEW_RELIC_LICENSE")
	log.Println(license)
	if license == "" {
		return
	}

	agent := gorelic.NewAgent()
	agent.NewrelicLicense = license
	agent.HTTPTimer = metrics.NewTimer()
	agent.CollectHTTPStat = true
	agent.NewrelicName = "site_txgruppi_com"
	agent.Run()

	m.Use(func(c martini.Context) {
		startTime := time.Now()
		c.Next()
		agent.HTTPTimer.UpdateSince(startTime)
	})
}

func registerLinksDao(m *martini.ClassicMartini, db *mgo.Database) {
	m.Map(links.NewDAO(db.C("links")))
}

func registerUrlShortener(m *martini.ClassicMartini, db *mgo.Database) {
	m.Map(urlshortener.New(db.C("urlshortener")))
}

func registerRenderer(m *martini.ClassicMartini) {
	m.Use(render.Renderer(render.Options{
		Layout: "",
	}))
}

func getRoot(r render.Render) {
	infoComment := "<!-- " + martini.Env + " -->"
	r.HTML(200, "main", map[string]interface{}{"env": template.HTML(infoComment)})
}

func getLinks(l *links.LinksDAO) []byte {
	links, err := l.All()
	panicIfErr(err)

	bytes, err := json.Marshal(links)
	panicIfErr(err)

	return bytes
}

func showHashidForId(params martini.Params, us *urlshortener.UrlShortener) string {
	id, err := strconv.Atoi(params["id"])
	if err == nil {
		hashid, err := us.HashIdFor(id)
		if err == nil {
			return hashid
		}
	}
	return "Can't generate hashid"
}

func redirectToUrl(res http.ResponseWriter, req *http.Request, r render.Render, params martini.Params, us *urlshortener.UrlShortener) {
	url, err := us.UrlFor(params["id"])
	if err == nil && url != "" {
		http.Redirect(res, req, url, 302)
	} else {
		r.HTML(404, "not_found", nil)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
