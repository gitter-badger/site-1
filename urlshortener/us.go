package urlshortener

import (
	"time"

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

	result, err := us.findByCriteria(c, bson.M{"_id": ids[0]})
	if err != nil {
		return "", nil
	}

	return result.Url, nil
}

func (us *UrlShortener) findUrlByAlias(c *mgo.Collection, alias string) (string, error) {
	result, err := us.findByCriteria(c, bson.M{"aliases": alias})
	if err != nil {
		return "", nil
	}

	return result.Url, nil
}

func (us *UrlShortener) findByCriteria(c *mgo.Collection, criteria bson.M) (ShortenedUrl, error) {
	result := ShortenedUrl{}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"lastUsed": time.Now()}},
		ReturnNew: true,
	}

	_, err := c.Find(criteria).Apply(change, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (us *UrlShortener) HashIdFor(id int) (string, error) {
	return us.hashid.Encrypt([]int{id})
}

func New(collection *mgo.Collection) *UrlShortener {
	h := hashids.New()
	h.Salt = "txgruppi" // It is not a secret

	return &UrlShortener{collection, h}
}
