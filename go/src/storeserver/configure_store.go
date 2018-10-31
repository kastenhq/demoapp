package storeserver

import (
	"crypto/tls"
	"models"
	"net/http"
	"os"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	graceful "github.com/tylerb/graceful"

	"store"
	"store/local"
	"storeserver/operations"
)

const (
	localHDDPAthEnvName = "LOCAL_STORE_PATH"
)

func configureFlags(api *operations.StoreAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.StoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	var storer store.Storer

	if lpath, ok := os.LookupEnv(localHDDPAthEnvName); ok {
		storer = &local.HDD{StorePath: lpath}
		log.Infof("Configured to use local HDD with path %s", lpath)
	} else {
		storer = &local.HDD{StorePath: os.TempDir()}
		log.Infof("Configured to use local HDD with TMP path %s", os.TempDir())
	}

	api.DeleteImageDataHandler = operations.DeleteImageDataHandlerFunc(func(params operations.DeleteImageDataParams) middleware.Responder {
		err := storer.Delete(params.HTTPRequest.Context(), params.ImageItem)
		if err != nil {
			log.Errorf("Failed to delete image " + err.Error())
			return operations.NewGetImageDataDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to delete image " + err.Error()})
		}
		return operations.NewDeleteImageDataOK()

	})
	api.GetImageDataHandler = operations.GetImageDataHandlerFunc(func(params operations.GetImageDataParams) middleware.Responder {
		imgData, err := storer.Read(params.HTTPRequest.Context(), params.ImageItem)
		if err != nil {
			log.Errorf("Failed to get image " + err.Error())
			return operations.NewGetImageDataDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to get image " + err.Error()})
		}
		return operations.NewGetImageDataOK().WithPayload(imgData)
	})

	api.StoreImageDataHandler = operations.StoreImageDataHandlerFunc(func(params operations.StoreImageDataParams) middleware.Responder {
		log.Infof("Store image image params %+v", params.ImageItem)
		img, err := storer.Write(params.HTTPRequest.Context(), params.ImageItem)
		if err != nil {
			log.Errorf("Failed to store image " + err.Error())
			return operations.NewGetImageDataDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to store image " + err.Error()})
		}
		return operations.NewStoreImageDataCreated().WithPayload(img)
	})

	healthzOK := operations.NewHealthzOK().WithPayload(&models.ServiceInfo{Version: "0.0.1"})
	api.HealthzHandler = operations.HealthzHandlerFunc(func(params operations.HealthzParams) middleware.Responder {
		return healthzOK
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
