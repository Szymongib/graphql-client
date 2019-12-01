package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/vrischmann/envconfig"

	"github.com/szymongib/graphql-client/test/schema"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Address  string `envconfig:"default=localhost:8001"`
	Endpoint string `envconfig:"default=/graphql"`
}

var (
	config   Config
	resolver *schema.Resolver

	apiAddress    string
	errorsAddress string
	nonGQLAddress string
)

func TestMain(m *testing.M) {
	err := envconfig.InitWithPrefix(&config, "APP")
	if err != nil {
		logrus.Errorf("Failed to read config: %s", err.Error())
		os.Exit(1)
	}

	code := setupTests(m)

	os.Exit(code)
}

func setupTests(m *testing.M) int {
	r := schema.NewResolver()
	resolver = &r

	gqlCfg := schema.Config{
		Resolvers: resolver,
	}
	executableSchema := schema.NewExecutableSchema(gqlCfg)

	log.Printf("Registering endpoint on %s...", config.Endpoint)

	router := mux.NewRouter()
	router.Use(func(i http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), schema.HeadersContextKey, r.Header))

			i.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/", handler.Playground("Dataloader", config.Endpoint))
	router.HandleFunc(config.Endpoint, handler.GraphQL(executableSchema))

	router.HandleFunc("/error/noGQL", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	})
	router.HandleFunc("/noGQL", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
		writer.WriteHeader(http.StatusOK)
	})

	http.Handle("/", router)

	go func() {
		log.Printf("API listening on %s...", config.Address)
		if err := http.ListenAndServe(config.Address, router); err != nil {
			panic(err)
		}
	}()

	apiAddress = fmt.Sprintf("http://%s%s", config.Address, config.Endpoint)
	errorsAddress = fmt.Sprintf("http://%s%s", config.Address, "/errors")
	nonGQLAddress = fmt.Sprintf("http://%s%s", config.Address, "/noGQL")

	time.Sleep(1 * time.Second)

	return m.Run()
}

func newResolver() *schema.Resolver {
	return &schema.Resolver{}
}
