// router.go
package router

import (
	"net/http"

	controller "github.com/Hifzu04/myMongoServer/Controller"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/movies", controller.GetAllMovies).Methods(http.MethodGet)
	api.HandleFunc("/movie", controller.InsertOneMovie).Methods(http.MethodPost)
	api.HandleFunc("/movie/{id}", controller.UpdateOneMovie).Methods(http.MethodPut)
	api.HandleFunc("/movie/{id}", controller.DeleteOneMovie).Methods(http.MethodDelete)
	api.HandleFunc("/movies", controller.DeleteAllMovies).Methods(http.MethodDelete)

	// Global CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(r)
}
