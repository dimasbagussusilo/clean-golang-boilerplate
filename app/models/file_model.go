package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type File struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileName     string             `bson:"file_name" json:"file_name" validate:"required"`
	OriginalName string             `bson:"original_name" json:"original_name" validate:"required"`
	FilePath     string             `bson:"file_path" json:"file_path" validate:"required"`
	FileSize     int64              `bson:"file_size" json:"file_size" validate:"required"`
	FileType     string             `bson:"file_type" json:"file_type" validate:"required"`
	CreatedAt    *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
