package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	"net/http"
	"os"
	"path"
	"time"
)

type SocialLink struct {
	Title   string
	Url     string
	SafeUrl template.URL
	Image   string
}

type SocialLinks []SocialLink

type Config struct {
	FileName string
	ModTime  time.Time
	Data     SocialLinks
	Debug    bool
}

func (c *Config) fileInfo() (os.FileInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path.Join(cwd, c.FileName))
	if err != nil {
		return nil, err
	}

	return fi, nil
}

func (c *Config) needReload(fi os.FileInfo) bool {
	if c.ModTime.Before(fi.ModTime()) {
		return true
	}
	return false
}

func (c *Config) Load() error {
	c.Log("Loading config")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.Open(path.Join(cwd, c.FileName))
	if err != nil {
		return err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	if err := dec.Decode(&c.Data); err != nil {
		return err
	}
	return nil
}

func (c *Config) Get() (SocialLinks, error) {
	fi, err := c.fileInfo()
	if err != nil {
		return nil, err
	}
	if c.needReload(fi) {
		err := c.Load()
		if err != nil {
			return nil, err
		}
	}
	c.ModTime = fi.ModTime()
	return c.Data, nil
}

func (c *Config) Log(msg string) {
	if !c.Debug {
		return
	}
	fmt.Println("[config] " + msg)
}

func main() {
	config := Config{"links.json", time.Time{}, nil, true}
	targetDate, _ := time.Parse("2006-01-02 15:04:05", "2014-05-03 20:00:00")

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		links, err := config.Get()
		if err != nil {
			panic(err)
		}

		for key := range links {
			if links[key].Title == "Skype" && links[key].Url[0:8] == "skype://" {
				links[key].SafeUrl = template.URL(links[key].Url)
			}
		}

		now := time.Now().Add(-(3 * time.Hour))
		daysLeft := int(targetDate.Sub(now).Hours() / 24)
		if (daysLeft < 0) {
			daysLeft = 0
		}

		r.HTML(200, "index", map[string]interface{}{"links": links, "daysLeft": daysLeft})
	})

	m.Get("/cs-logica", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "http://code-squad.com/curso/logica-programacao/avulso", 302)
	})

	m.Get("/sn-php", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "http://www.schoolofnet.com/cursos/php-basico/", 302)
	})

	m.Get("/pti-50-gratis", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "http://www.profissionaisti.com.br/2013/03/os-50-melhores-cursos-gratis-de-ti-de-toda-a-internet/", 302)
	})

	m.Get("/rls-java-gratis", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "http://www.rlsystem.com.br/curso-java-gratis/", 302)
	})

	m.Get("/rls-android-gratis", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "http://www.rlsystem.com.br/curso-android-gratis/", 302)
	})

	m.Get("/php-velha", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "https://dl.dropboxusercontent.com/u/1274888/php-velha-v0.zip", 302)
	})

	m.NotFound(func(r render.Render) {
		r.HTML(200, "not_found", nil)
	})

	m.Run()
}
