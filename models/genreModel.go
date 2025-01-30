package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	Id         bson.ObjectID `bson:"_id"`
	Name       *string       `json:"name" validate:"required,min=4,max=100"`
	Created_at time.Time     `json:"created_at"`
	Updated_at time.Time     `json:"updated_at"`
	Genre_id   int           `json:"genre_id" validate:"required"`
}
