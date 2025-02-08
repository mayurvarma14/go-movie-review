package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id"`
	Name         *string       `json:"name" validate:"required,min=4,max=100"`
	Username     *string       `json:"username" validate:"required,min=4,max=100"`
	Password     *string       `json:"password" validate:"required,min=8"`
	Email        *string       `json:"email" validate:"email,required"`
	Token        *string       `json:"token,omitempty" bson:"token,omitempty"`
	UserType     *string       `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string       `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
	UserID       string        `json:"user_id" bson:"user_id"`
}
