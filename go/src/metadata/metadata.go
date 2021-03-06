package metadata

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"models"
	"restclient"
	"transport"
)

const (
	// StoreServiceAddrEnv allows to specify store service
	StoreServiceAddrEnv = "STORE-IP"
	port                = "8000"
)

// MetaDater interface to deal with metadata
// English language: metadata + er = thing which handels metadata
type MetaDater interface {
	GetAllImages(context.Context) (models.ImageList, error)
	FetchImage(context.Context, strfmt.UUID) (models.ImageData, error)
	FindImages(context.Context, map[string]interface{}) (models.ImageList, error)
	Delete(context.Context, *models.ImageMeta) error
	Add(context.Context, models.ImageData) (*models.Image, error)
}

// StoreClient discover and returns StoreRest client
func StoreClient() *restclient.Rest {
	storeAddr := "store"
	if tmpAddr, ok := os.LookupEnv(StoreServiceAddrEnv); ok {
		storeAddr = tmpAddr
	}

	host := fmt.Sprintf("%s:%s", storeAddr, port)
	cfg := restclient.DefaultTransportConfig().WithHost(host)
	httpClient := &http.Client{Transport: transport.NewTracingTransport(http.DefaultTransport)}
	ct := client.NewWithClient(cfg.Host, cfg.BasePath, cfg.Schemes, httpClient)
	return restclient.New(ct, strfmt.Default)
}
