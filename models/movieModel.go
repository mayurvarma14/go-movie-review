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
	MovieID   int           `json:"movie_id" bson:"movie_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
