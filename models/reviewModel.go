package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Reviews struct {
	ID         bson.ObjectID `bson:"_id"`
	MovieID    int           `json:"movie_id" bson:"movie_id"`
	ReviewerID bson.ObjectID `json:"reviewer_id" bson:"reviewer_id"`
	Review     *string       `json:"review" validate:"required"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}
