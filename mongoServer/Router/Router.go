package router

import (
	controller "github.com/Hifzu04/myMongoServer/Controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	myroute := mux.NewRouter()

	myroute.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	myroute.HandleFunc("/api/movie", controller.InsertOneMovie).Methods("POST")
	myroute.HandleFunc("/api/movie/{id}", controller.UpdateOneMovie).Methods("PUT")
	myroute.HandleFunc("/api/movie/delete/{id}", controller.DeleteOneMovie).Methods("DELETE")
	myroute.HandleFunc("/api/movies/delete", controller.DeleteAllMovies).Methods("DELETE")
	return myroute
}
