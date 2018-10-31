package mongo

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	. "gopkg.in/check.v1"

	"metadata"
	"models"
	"uuid"
)

const (
	testImage     = "testimage/logo.png"
	testImageName = "logo.png"
)

type MongoMetaDataSuite struct {
	db        Mongo
	mongoDRes *dockertest.Resource
	storeDRes *dockertest.Resource
	testImage models.Image
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MongoMetaDataSuite{})

func (s *MongoMetaDataSuite) SetUpSuite(c *C) {
	fbuf, err := ioutil.ReadFile(testImage)
	c.Assert(err, IsNil)
	c.Assert(fbuf, NotNil)

	s.testImage = models.Image{
		Meta: &models.ImageMeta{
			CreationTime: models.TimeStamp(time.Now()),
			ID:           uuid.New(),
		},
		Base64: models.ImageData(base64.StdEncoding.EncodeToString(fbuf)),
	}

	dOpts := &dockertest.RunOptions{
		Repository:   "bitnami/mongodb",
		Tag:          "3.6",
		Env:          []string{"MONGODB_DATABASE=images", "MONGODB_USERNAME=testuser", "MONGODB_PASSWORD=testpassword"},
		PortBindings: map[dc.Port][]dc.PortBinding{"27017/tcp": {{HostPort: "27017"}}},
	}

	pool, err := dockertest.NewPool("")
	c.Assert(err, IsNil)
	s.mongoDRes, err = pool.RunWithOptions(dOpts)
	c.Assert(err, IsNil)
	if err := pool.Retry(func() error {
		_, err := net.DialTimeout("tcp", net.JoinHostPort("", "27017"), time.Second*5)
		if err != nil {
			c.Logf("Failed to dial to mongo: ", err)
			return err
		}
		return nil
	}); err != nil {
		c.Errorf("Failed to connect to mongo %s", err.Error())
	}

	s.db = Mongo{DBurl: "mongodb://testuser:testpassword@localhost:27017/?authSource=images&authMechanism=SCRAM-SHA-1"}
	s.startStore(c)
	os.Setenv(metadata.StoreServiceAddrEnv, "localhost")
	// bitnami mongo is slow :(
	time.Sleep(time.Second * 10)
}

func (s *MongoMetaDataSuite) TearDownSuite(c *C) {
	s.db.Conn.Close()
	err := s.mongoDRes.Close()
	c.Assert(err, IsNil)
	err = s.storeDRes.Close()
	c.Assert(err, IsNil)
}

func (s *MongoMetaDataSuite) TestPing(c *C) {
	err := s.db.Ping()
	c.Assert(err, IsNil)
	serv := s.db.Conn.LiveServers()
	for _, s := range serv {
		c.Logf("%s\n", s)
	}
}

func (s *MongoMetaDataSuite) TestAdd(c *C) {
	err := s.db.Ping()
	c.Assert(err, IsNil)
	img, err := s.db.Add(context.Background(), s.testImage.Base64)
	c.Assert(err, IsNil)
	c.Assert(*img, FitsTypeOf, models.Image{})
}

func (s *MongoMetaDataSuite) TestGetAllImages(c *C) {
	err := s.db.Ping()
	c.Assert(err, IsNil)
	_, err = s.db.Add(context.Background(), s.testImage.Base64)
	imgs, err := s.db.GetAllImages(context.Background())
	c.Assert(err, IsNil)
	for _, img := range imgs {
		c.Logf("List: %+v", *img)
	}
	c.Assert(len(imgs) > 0, Equals, true)
}

func (s *MongoMetaDataSuite) TestDelete(c *C) {
	err := s.db.Ping()
	c.Assert(err, IsNil)
	_, err = s.db.Add(context.Background(), s.testImage.Base64)
	imgs, err := s.db.GetAllImages(context.Background())
	c.Assert(err, IsNil)
	err = s.db.Delete(context.Background(), imgs[0])
}

func (s *MongoMetaDataSuite) startStore(c *C) {
	dOpts := &dockertest.RunOptions{
		Repository:   "store-server",
		Tag:          "latest",
		PortBindings: map[dc.Port][]dc.PortBinding{"8000/tcp": {{HostPort: "8000"}}},
	}

	pool, err := dockertest.NewPool("")
	c.Assert(err, IsNil)

	s.storeDRes, err = pool.RunWithOptions(dOpts)
	c.Assert(err, IsNil)
	if err := pool.Retry(func() error {
		_, err := net.DialTimeout("tcp", net.JoinHostPort("", "8000"), time.Second*5)
		if err != nil {
			c.Logf("Failed to dial store server:", err)
			return err
		}
		return nil
	}); err != nil {
		c.Errorf("Failed to connect to store %s", err.Error())
	}
}
