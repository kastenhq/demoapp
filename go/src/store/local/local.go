package local

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"models"
	"store"
	"uuid"
)

var _ store.Storer = (*HDD)(nil)

// HDD structre for HDD storer interface
type HDD struct {
	StorePath string
}

// Write writes Image into local storage
func (s *HDD) Write(ctx context.Context, image *models.Image) (*models.Image, error) {
	err := os.Mkdir(s.StorePath, os.ModeDir)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	if image.Meta.ID == "" {
		image.Meta.ID = uuid.New()
		image.Meta.CreationTime = models.TimeStamp(time.Now())
	}
	image.Meta.Location = filepath.Join(s.StorePath, string(image.Meta.ID))
	return image, ioutil.WriteFile(image.Meta.Location, []byte(image.Base64), 0644)
}

// Read reads local data for provided ImageMeta
func (s *HDD) Read(ctx context.Context, image *models.ImageMeta) (models.ImageData, error) {
	fb, err := ioutil.ReadFile(image.Location)
	if err != nil {
		return "", err
	}
	return models.ImageData(string(fb)), nil
}

// Delete deletes local data for provided image meta
func (s *HDD) Delete(ctx context.Context, image *models.ImageMeta) error {
	return os.Remove(image.Location)
}
