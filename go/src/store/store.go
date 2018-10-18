package store

import (
	"context"

	"models"
)

// Storer interface to implement different types of storages
type Storer interface {
	Write(context.Context, *models.Image) (*models.Image, error)
	Read(context.Context, *models.ImageMeta) (models.ImageData, error)
	Delete(context.Context, *models.ImageMeta) error
}
