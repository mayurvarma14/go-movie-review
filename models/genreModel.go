package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      *string       `json:"name" validate:"required,min=4,max=100"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	GenreID   int           `json:"genre_id" validate:"required" bson:"genre_id"`
}
