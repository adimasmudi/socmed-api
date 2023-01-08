package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	Id      primitive.ObjectID `json:"id,omitempty"`
	Owner   string             `json:"owner,omitempty" validate:"required"`
	Title   string             `json:"title,omitempty" validate:"required"`
	Content string             `json:"content,omitempty" validate:"required"`
}
