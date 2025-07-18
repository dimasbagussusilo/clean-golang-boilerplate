package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// StoreLocation represents the location details of a store
type StoreLocation struct {
	Address    string  `bson:"address" json:"address" validate:"required"`
	City       string  `bson:"city" json:"city" validate:"required"`
	PostalCode string  `bson:"postal_code" json:"postal_code"`
	Country    string  `bson:"country" json:"country" validate:"required"`
	Latitude   float64 `bson:"latitude" json:"latitude"`
	Longitude  float64 `bson:"longitude" json:"longitude"`
}

type Store struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Type        string             `bson:"type" json:"type" validate:"required,oneof=retail grosir"`
	//IsActive    bool               `bson:"is_active" json:"is_active" validate:"oneof=true false" default:"true"`
	Location  *StoreLocation `bson:"location" json:"location"`
	CreatedAt *time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
