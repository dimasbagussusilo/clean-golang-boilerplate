package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Store struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Type        string             `bson:"type" json:"type" validate:"required,oneof=retail grosir"`
	//IsActive    bool               `bson:"is_active" json:"is_active" validate:"oneof=true false" default:"true"`
	CreatedAt *time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
