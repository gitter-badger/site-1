package db

import (
	"errors"
	"labix.org/v2/mgo"
	"strings"
)

func NewMongoDB(url string) MongoDB {
	if !strings.HasPrefix(url, "mongodb://") {
		url = "mongodb://" + url
	}
	if !strings.Contains(url[10:], "/") {
		url = url + "/"
	}
	databaseName := url[strings.LastIndex(url, "/")+1:]
	return MongoDB{url, databaseName, nil}
}

type MongoDB struct {
	url          string
	databaseName string
	session      *mgo.Session
}

func (m *MongoDB) Dial() (err error) {
	m.session, err = mgo.Dial(m.url)
	return
}

func (m *MongoDB) Close() {
	if m.session != nil {
		m.session.Close()
	}
}

func (m *MongoDB) Database() (db *mgo.Database, err error) {
	if m.databaseName == "" {
		err = errors.New(DATABASE_NAME_EMPTY)
		return
	}
	if m.session != nil {
		db = m.session.DB(m.databaseName)
	}
	return
}

func (m *MongoDB) Url() string {
	return m.url
}
