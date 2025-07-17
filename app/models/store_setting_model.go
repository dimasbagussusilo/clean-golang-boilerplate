package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StoreSetting struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StoreID   primitive.ObjectID `bson:"store_id" json:"store_id"`
	IsHaram   bool               `bson:"is_haram" json:"is_haram"`
	CreatedAt *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
