package links

import (
	"errors"

	"labix.org/v2/mgo"
)

func NewDAO(collection *mgo.Collection) *LinksDAO {
	return &LinksDAO{collection}
}

type LinksDAO struct {
	collection *mgo.Collection
}

func (d *LinksDAO) All() (Links, error) {
	result := Links{}

	if d.collection == nil {
		return result, errors.New(NIL_COLLECTION)
	}

	err := d.collection.Find(nil).Sort("order").All(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
