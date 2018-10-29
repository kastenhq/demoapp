package mongo

import (
	"context"
	"errors"

	"github.com/globalsign/mgo"

	"metadata"
	"models"
	"restclient/operations"
)

const (
	dbName   = "images"
	collName = "metadata"
)

var _ metadata.MetaDater = (*Mongo)(nil)

// Mongo structre for Mongo metadater interface
type Mongo struct {
	DBurl string
	Conn  *mgo.Session
}

// GetAllImages returns all images meta data
func (s *Mongo) GetAllImages(ctx context.Context) (models.ImageList, error) {
	err := s.Ping()
	if err != nil {
		return models.ImageList{}, err
	}
	c := s.Conn.DB(dbName).C(collName)
	imgs := models.ImageList{}
	return imgs, c.Find(nil).All(&imgs)
}

// FindImages performs a select and returns list of images
func (s *Mongo) FindImages(ctx context.Context, tags map[string]string) (models.ImageList, error) {
	return models.ImageList{}, nil
}

// Delete removes meta data for porvided imagemeta
func (s *Mongo) Delete(ctx context.Context, imgMeta *models.ImageMeta) error {
	err := s.Ping()
	if err != nil {
		return err
	}
	c := s.Conn.DB(dbName).C(collName)

	return c.Remove(*imgMeta)
}

// Add call store service and creates meta data
func (s *Mongo) Add(ctx context.Context, imgData models.ImageData) (*models.Image, error) {
	err := s.Ping()
	if err != nil {
		return nil, err
	}

	img := &models.Image{
		Base64: imgData,
		Meta:   &models.ImageMeta{},
	}

	store := metadata.StoreClient()
	addImg := &operations.StoreImageDataParams{
		Context:   ctx,
		ImageItem: img,
	}

	out, err := store.Operations.StoreImageData(addImg)
	if err != nil {
		return nil, err
	}

	img = out.Payload
	c := s.Conn.DB(dbName).C(collName)
	return img, c.Insert(*img.Meta)
}

// Ping ping current session or creates a new one
func (s *Mongo) Ping() error {
	var err error
	if s.DBurl == "" {
		return errors.New("DBurl is empty")
	}
	if s.Conn == nil {
		s.Conn, err = mgo.Dial(s.DBurl)
		return err
	}
	return s.Conn.Ping()
}
