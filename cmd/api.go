package main

import (
	"log"
	"net/http"
	"time"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
	"github.com/CassRamos/go-ecom-api.git/internal/orders"
	"github.com/CassRamos/go-ecom-api.git/internal/products"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // Assigns a unique ID to each request and stores it in the request context.
	r.Use(middleware.RealIP)    // Sets the RemoteAddr to the value of X-Forwarded-For or X-Real-IP headers, if present.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal to cancel the request if it takes longer than 60 seconds to complete.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check passed."))
	})

	productService := products.NewService(repository.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.GetProductById)
	r.Post("/products", productHandler.CreateProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	orderService := orders.NewService(repository.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService)
	r.Post("/orders", orderHandler.PlaceOrder)
	r.Get("/orders/{id}", orderHandler.GetOrderByID)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)
	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
