package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/zippy/config"
	_ "github.com/rosset7i/zippy/docs"
	"github.com/rosset7i/zippy/internal/infra/database"
	"github.com/rosset7i/zippy/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Zippy API
// @version 1.0
// @description IDK.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config := config.LoadConfig()

	db, err := database.NewPool(context.Background(), config.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userHandler := handlers.NewUserHandler(database.NewUser(db), config)
	productHandler := handlers.NewProductHandler(database.NewProduct(db))

	routes := chi.NewRouter()
	routes.Use(middleware.Logger)
	routes.Use(middleware.Recoverer)
	routes.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})
	routes.Route("/products", func(r chi.Router) {
		// r.Use(jwtauth.Verifier(config.TokenAuth))
		// r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.FetchPaged)
		r.Get("/{id}", productHandler.FetchById)
		r.Post("/", productHandler.Create)
		r.Put("/", productHandler.Update)
		r.Delete("/", productHandler.Delete)
	})
	routes.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:7000/docs/doc.json")))
	if err := http.ListenAndServe(config.WebServerAddress, routes); err != nil {
		log.Fatal(err)
	}
}
