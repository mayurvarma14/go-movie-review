package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      *string       `json:"name" validate:"required"`
	Topic     *string       `json:"topic" validate:"required"`
	GenreID   int           `json:"genre_id" bson:"genre_id"`
	MovieURL  *string       `json:"movie_url" validate:"required"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
