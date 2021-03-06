package metaserver

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	graceful "github.com/tylerb/graceful"

	"metadata"
	"metadata/mongo"
	"metaserver/operations"
	"models"
	"requestid"
	"tracing"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name meta --spec ../../../swagger.yaml --server-package metaserver --client-package restclient --operation healthz --operation addImage --operation listImages --operation getImage --operation deleteImage --skip-models

func configureFlags(api *operations.MetaAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MetaAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.ServerShutdown = func() {}
	log.Info("Configuring tracing")
	if closer, err := tracing.Configure("metadata"); err != nil {
		log.Errorf("Unable to configure tracing. %s", err.Error())
	} else {
		tmp := api.ServerShutdown
		api.ServerShutdown = func() {
			closer.Close()
			tmp()
		}
	}

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	var metadata metadata.MetaDater

	if mHost, ok := os.LookupEnv(mongo.HOSTENV); ok {
		mPass := os.Getenv(mongo.PASSWORDENV)
		mUser := os.Getenv(mongo.USERNAMEENV)
		mDBUrl := fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=images&authMechanism=SCRAM-SHA-1", mUser, mPass, mHost)
		metadata = &mongo.Mongo{DBurl: mDBUrl}
	} else {
		log.Panic("Unsupported metadata provider type")
	}

	api.AddImageHandler = operations.AddImageHandlerFunc(func(params operations.AddImageParams) middleware.Responder {
		rImg, err := metadata.Add(params.HTTPRequest.Context(), params.ImageItem.Base64)
		log.Info("Add image params ")
		if err != nil {
			log.Errorf("Failed to add new Image meatdata with err %s", err.Error())
			return operations.NewAddImageDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to add new image " + err.Error()})
		}
		return operations.NewAddImageCreated().WithPayload(rImg)
	})

	api.DeleteImageHandler = operations.DeleteImageHandlerFunc(func(params operations.DeleteImageParams) middleware.Responder {
		err := metadata.Delete(params.HTTPRequest.Context(), &models.ImageMeta{ID: params.ItemID})
		if err != nil {
			log.Errorf("Failed to delete Image with err %s", err.Error())
			return operations.NewDeleteImageDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to delete image " + err.Error()})
		}
		return operations.NewDeleteImageOK()
	})

	api.GetImageHandler = operations.GetImageHandlerFunc(func(params operations.GetImageParams) middleware.Responder {
		rImg, err := metadata.FetchImage(params.HTTPRequest.Context(), params.ItemID)
		if err != nil {
			log.Errorf("Failed to fetch Image with err %s", err.Error())
			return operations.NewGetImageDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to fetch image " + err.Error()})
		}
		ret, _ := base64.StdEncoding.DecodeString(string(rImg))
		return operations.NewGetImageOK().WithPayload(ret)
	})

	healthzOK := operations.NewHealthzOK().WithPayload(&models.ServiceInfo{Version: "0.0.1"})
	api.HealthzHandler = operations.HealthzHandlerFunc(func(params operations.HealthzParams) middleware.Responder {
		return healthzOK
	})

	api.ListImagesHandler = operations.ListImagesHandlerFunc(func(params operations.ListImagesParams) middleware.Responder {
		imgs, err := metadata.GetAllImages(params.HTTPRequest.Context())
		if err != nil {
			log.Errorf("Failed to list images %s", err)
			return operations.NewListImagesDefault(500).WithPayload(&models.ErrorDetail{Message: "Failed to list images " + err.Error()})
		}
		return operations.NewListImagesOK().WithPayload(imgs)
	})

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
	return tracing.Middleware(requestid.PropagateOrGenerate(handler))
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	// So this is a good place to plug in a panic handling middleware, logging and metrics
	return handler
}
