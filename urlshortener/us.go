package urlshortener

import (
	"github.com/speps/go-hashids"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type ShortenedUrl struct {
	Id  int `bson:"_id"`
	Url string
}

type UrlShortener struct {
	collection *mgo.Collection
	hashid     *hashids.HashID
}

func (us *UrlShortener) UrlFor(id string) (url string, err error) {
	defer func() {
		recover()
	}()

	url, err = us.findUrlByAlias(us.collection, id)
	if url == "" {
		url, err = us.findUrlById(us.collection, id)
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

func (us *UrlShortener) HashIdFor(id int) (string, error) {
	return us.hashid.Encrypt([]int{id})
}

func New(collection *mgo.Collection) UrlShortener {
	h := hashids.New()
	h.Salt = "txgruppi" // It is not a secret

	us := UrlShortener{collection, h}

	return us
}
