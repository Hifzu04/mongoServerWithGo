package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/Hifzu04/myMongoServer/Router"
)

func main() {
	

	// uri := os.Getenv("MONGODB_URI")
	// if uri == "" {
	// 	fmt.Println("MONGO_URI not set")
	// 	return
	// }
	// fmt.Println("MONGO_URI is set to", uri)

	fmt.Println("welcome to mongodb server")

	r := router.Router()

	log.Fatal(http.ListenAndServe(":27017", r))

	fmt.Println("listening and serving")

}
