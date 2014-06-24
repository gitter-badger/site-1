package links

import (
	"labix.org/v2/mgo/bson"
)

type Link struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Title string        `bson:"title" json:"title"`
	Url   string        `bson:"url" json:"url"`
	Image string        `bson:"image" json:"image"`
	Order int           `bson:"order" json:"order"`
}

type Links []Link
