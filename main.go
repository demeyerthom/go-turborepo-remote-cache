package main

import (
	"context"
	"github.com/demeyerthom/go-turborepo-remote-cache/internal"
	"github.com/demeyerthom/go-turborepo-remote-cache/internal/storage"
	"github.com/getkin/kin-openapi/openapi3filter"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package=internal -generate=gorilla-server,types,spec -o=internal/server.generated.go https://turbo.build/api/remote-cache-spec

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{
		PadLevelText: true,
	})

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func buildHandler() (http.Handler, error) {
	// Pick type of storage from env
	s, err := storage.CreateStorage(storage.GCS)
	if err != nil {
		return nil, err
	}

	si := internal.NewServerHandler(s)

	specs, err := internal.GetSwagger()
	if err != nil {
		return nil, err
	}
	specs.Servers = nil

	// TODO: add simple authentication middleware
	return internal.HandlerWithOptions(si, internal.GorillaServerOptions{
		Middlewares: []internal.MiddlewareFunc{
			nethttpmiddleware.OapiRequestValidatorWithOptions(specs, &nethttpmiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
						return nil
					},
				},
			}),
		},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}), nil
}

func main() {
	h, err := buildHandler()
	if err != nil {
		log.Fatalf("could not build handler: %v", err)
	}

	s := &http.Server{
		Handler: h,
		//TODO make port configurable via env
		Addr: ":3000",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(func() error {
		log.Printf("Server is running on %s", s.Addr)
		return s.ListenAndServe()
	}())
}
