package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Reviews struct {
	Id          bson.ObjectID `bson:"_id"`
	Movie_id    int           `json:"movie_id"`
	Reviewer_id string        `json:"reviewer_id"`
	Review      *string       `json:"review" validate:"required"`
	Created_at  time.Time     `json:"created_at"`
	Updated_at  time.Time     `json:"updated_at"`
}
