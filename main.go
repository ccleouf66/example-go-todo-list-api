package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codenito/example-go-todo-list-api/pkg/api"
	"github.com/codenito/example-go-todo-list-api/pkg/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	mongoConf := &store.StoreOption{
		Address:    os.Getenv("MONGO_ADDR"),
		User:       os.Getenv("MONGO_USR"),
		Password:   os.Getenv("MONGO_PWD"),
		DbName:     os.Getenv("MONGO_DB"),
		RsName:     os.Getenv("MONGO_RS"),
		AuthSource: os.Getenv("MONGO_AUTH_SOURCE"),
	}

	ctx := context.Background()

	db, err := store.NewMongoStore(ctx, mongoConf)
	if err != nil {
		log.Fatalln(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "application/json"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	taskHandler := api.TaskHandler{
		Store: db,
	}

	metricsHandler := api.NewLetricsHandler()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(metricsHandler.IncrementTotalQueryMetric)
		r.Route("/task", taskHandler.ServeHTTP)
	})

	r.Get("/metrics", metricsHandler.GetMetrics())

	// Update fake metric
	go func() {
		for {
			metricsHandler.RandDesiredPodNumber()
			time.Sleep(30 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", r))
}
