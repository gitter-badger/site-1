package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
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

		r.HTML(200, "index", links)
	})

	m.Run()
}
