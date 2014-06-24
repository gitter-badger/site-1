package links

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo"
	"os"
	"testing"
)

func TestThrowWithNilCollection(t *testing.T) {
	d := NewDAO(nil)
	_, err := d.All()
	assert.Equal(t, err.Error(), NIL_COLLECTION)
}

func TestFindAll(t *testing.T) {
	_, _, c := getMongoObjects()
	clearCollection(c)
	insertFakeData(c)

	dao := NewDAO(c)

	expected := Links{}
	err := c.Find(nil).Sort("order").All(&expected)
	if err != nil {
		panic(err)
	}

	actual, err := dao.All()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(expected), len(actual))

	assert.True(t, assert.ObjectsAreEqual(expected, actual))

	assert.Equal(t, 0, actual[0].Order)
	assert.Equal(t, 1, actual[1].Order)
	assert.Equal(t, 2, actual[2].Order)
}

func getMongoObjects() (*mgo.Session, *mgo.Database, *mgo.Collection) {
	s := connectToMongoDB()
	d := getTestDatabase(s)
	c := getTestCollection(d)
	return s, d, c
}

func connectToMongoDB() *mgo.Session {
	mongoUrl := os.Getenv("MONGO_URL")

	if mongoUrl == "" {
		return nil
	}

	s, _ := mgo.Dial(mongoUrl)

	return s
}

func getTestDatabase(s *mgo.Session) *mgo.Database {
	return s.DB("com_txgruppi_site_test")
}

func getTestCollection(d *mgo.Database) *mgo.Collection {
	return d.C("links")
}

func clearCollection(c *mgo.Collection) {
	c.DropCollection()
}

func insertFakeData(c *mgo.Collection) {
	c.Insert(Link{
		"",
		"First",
		"http://first.test",
		"http://first.test/image",
		0,
	})
	c.Insert(Link{
		"",
		"Third",
		"http://third.test",
		"http://third.test/image",
		2,
	})
	c.Insert(Link{
		"",
		"Second",
		"http://second.test",
		"http://second.test/image",
		1,
	})
}
