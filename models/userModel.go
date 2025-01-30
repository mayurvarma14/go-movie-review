package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID `bson:"_id"`
	Name          *string       `json:"name" validate:"required,min=4,max=100"`
	Username      *string       `json:"username" validate:"required,min=4,max=100"`
	Password      *string       `json:"password" validate:"required,min=8"`
	Email         *string       `json:"email" validate:"email,required"`
	Token         *string       `json:"token"`
	User_type     *string       `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token *string       `json:"refresh_token"`
	Created_at    time.Time     `json:"created_at"`
	Updated_at    time.Time     `json:"updated_at"`
	User_id       string        `json:"user_id"`
}
