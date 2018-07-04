package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Todo struct {
	Id          bson.ObjectId `bson: "_id" json:"id, omitempty"`
	Title       string        `bson:"title" json:"title"`
	Contents    string        `bson:"contents" json:"contents"`
	Location    string        `bson:"location" json:"location"`
	Deadline    time.Time     `bson:"deadline" json:"deadline"`
	CreatedTime time.Time     `bson:"created_timestamp" json:"created_timestamp"`
}
