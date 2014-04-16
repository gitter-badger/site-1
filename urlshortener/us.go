package urlshortener

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"github.com/speps/go-hashids"
	"strings"
)

type ShortenedUrl struct {
	Id int `bson:"_id"`
	Url string
}

type UrlShortener struct {
	mongoUrl string
	dbName string
	collectionName string
	hashid *hashids.HashID
}

func (us *UrlShortener) UrlFor(id string) (url string, err error) {
	defer func(){
		recover()
	}()

	session, err := mgo.Dial(us.mongoUrl)
	if err != nil {
		return "", err
	}
	defer session.Close()

	c := session.DB(us.dbName).C(us.collectionName)

	url, err = us.findUrlByAlias(c, id)
	if url == "" {
		url, err = us.findUrlById(c, id)
	}

	return
}

func (us *UrlShortener) findUrlById(c *mgo.Collection, id string) (string, error) {
	ids := us.hashid.Decrypt(id)
	if len(ids) != 1 {
		return "", nil
	}

	result := ShortenedUrl{}
	err := c.Find(bson.M{"_id": ids[0]}).One(&result)
	if err != nil {
		return "", err
	}

	return result.Url, nil
}

func (us *UrlShortener) findUrlByAlias(c *mgo.Collection, alias string) (string, error) {
	result := ShortenedUrl{}
	err := c.Find(bson.M{"aliases": alias}).One(&result)
	if err != nil {
		return "", err
	}

	return result.Url, nil
}

func New(collectionName string) *UrlShortener {
	mongoUrl := os.Getenv("MONGO_URL")

	if mongoUrl == "" {
		return nil
	}

	h := hashids.New()
	h.Salt = "txgruppi" // It is not a secret

	us := UrlShortener{mongoUrl, mongoUrl[strings.LastIndex(mongoUrl, "/")+1:], collectionName, h}

	return &us
}
