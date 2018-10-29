package local

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"models"
	"store"
	"uuid"
)

const (
	testImage     = "testimage/logo.png"
	testImageName = "logo.png"
)

type LocalStoreSuite struct {
	testImage  models.Image
	storer     store.Storer
	testFolder string
}

var _ = Suite(&LocalStoreSuite{})

func (s *LocalStoreSuite) SetUpSuite(c *C) {
	s.testFolder = c.MkDir()
	s.storer = &HDD{StorePath: s.testFolder}
	fbuf, err := ioutil.ReadFile(testImage)
	c.Assert(err, IsNil)
	c.Assert(fbuf, NotNil)

	s.testImage = models.Image{
		Meta: &models.ImageMeta{
			CreationTime: models.TimeStamp(time.Now()),
			ID:           uuid.New(),
			Location:     filepath.Join(s.testFolder, testImageName),
		},
		Base64: models.ImageData(base64.StdEncoding.EncodeToString(fbuf)),
	}
}

func Test(t *testing.T) { TestingT(t) }

func (s *LocalStoreSuite) TearDownSuite(c *C) {
	err := os.Remove(s.testFolder)
	c.Assert(err, IsNil)
}

func (s *LocalStoreSuite) TestALocalWrite(c *C) {
	img, err := s.storer.Write(context.TODO(), &s.testImage)
	c.Assert(err, IsNil)
	c.Assert(img.Meta.Location, Equals, s.testImage.Meta.Location)
}
func (s *LocalStoreSuite) TestBLocalRead(c *C) {
	img, err := s.storer.Read(context.TODO(), s.testImage.Meta)
	c.Assert(err, IsNil)
	c.Assert(img, Equals, s.testImage.Base64)
}

func (s *LocalStoreSuite) TestCLocalDelete(c *C) {
	err := s.storer.Delete(context.TODO(), s.testImage.Meta)
	c.Assert(err, IsNil)
	_, err = ioutil.ReadFile(s.testImage.Meta.Location)
	c.Assert(os.IsNotExist(err), Equals, true)
}
