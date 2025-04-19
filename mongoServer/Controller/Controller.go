package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "github.com/Hifzu04/myMongoServer/Model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionstring = ""
const dbname = "netflix"
const colname = "watchlist"

var collection *mongo.Collection

// connect with mongodb
func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionstring)

	//connet to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongodb connection sucess")

	collection = client.Database(dbname).Collection(colname)

	// if collection instance is ready ,
	fmt.Println("collection instance is ready")

}

// insert a data into mongodb in golang
// bring data from req body(url params) take it and insert the data into data base
// mongodb helper-file
// insert 1 record
func insertOneMovie(movie model.Netflix) { //take movie in model(stuct) format
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted oneMovie in db with id : ", inserted.InsertedID)
}

// update a record in mongodb in golang

// update 1 movie
func updateOneMovie(movieID string) {
	//now we have to convert movieID (string) to something mongoDB can understand
	id, _ := primitive.ObjectIDFromHex(movieID) //convert

	//filter that particular movie which needs to be updatead on the basis of id
	filter := bson.M{"_id": id} //key value
	update := bson.M{"$ set": bson.M{"Watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified course is : ", result.ModifiedCount)
}

// delete one and delete many movies in mongodb
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("movie got deleted with delete count : ", deleteCount)
}

func deleteAllMovie() int64 {
	dltres, _ := collection.DeleteMany(context.Background(), bson.D{{}})
	fmt.Println("deleted items , dlt count is : ", dltres.DeletedCount)
	return dltres.DeletedCount
}

// get all the movie , #tricky
func getAllMovies() []primitive.M {

	//it will return a cursor
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M
	//while we are getting the next cursor keep on looping
	for cursor.Next(context.Background()) {
		var movie bson.M

		err := cursor.Decode(&movie)

		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())

	return movies
}

//controller for all the above helpers

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func InsertOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix

	json.NewDecoder(r.Body).Decode(&movie)
	//send to helper
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func UpdateOneMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
