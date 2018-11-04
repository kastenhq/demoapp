package mongo

import (
	"context"
	"errors"
	"log"

	"github.com/globalsign/mgo"
	"github.com/go-openapi/strfmt"

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

// FetchImage searches thru metadata and fetchs image data from store
func (s *Mongo) FetchImage(ctx context.Context, id strfmt.UUID) (models.ImageData, error) {
	imgs, err := s.FindImages(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return "", err
	}
	log.Printf("Images %+v", imgs)
	if len(imgs) != 1 {
		return "", errors.New("Found multiple or none images for the same ID")
	}
	store := metadata.StoreClient()
	getImg := &operations.GetImageDataParams{
		Context:   ctx,
		ImageItem: imgs[0],
	}
	imgData, err := store.Operations.GetImageData(getImg)
	if err != nil {
		return "", err
	}
	return imgData.Payload, nil
}

// FindImages performs a select and returns list of images
func (s *Mongo) FindImages(ctx context.Context, tags map[string]interface{}) (models.ImageList, error) {
	err := s.Ping()
	if err != nil {
		return models.ImageList{}, err
	}
	c := s.Conn.DB(dbName).C(collName)
	imgs := models.ImageList{}
	err = c.Find(tags).All(&imgs)
	log.Printf("Found for search %+v", imgs)
	return imgs, err
}

// Delete removes meta data for porvided imagemeta
func (s *Mongo) Delete(ctx context.Context, imgMeta *models.ImageMeta) error {
	err := s.Ping()
	if err != nil {
		return err
	}
	c := s.Conn.DB(dbName).C(collName)

	store := metadata.StoreClient()
	delImg := &operations.DeleteImageDataParams{
		Context:   ctx,
		ImageItem: imgMeta,
	}
	_, err = store.Operations.DeleteImageData(delImg)
	if err != nil {
		return err
	}

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
