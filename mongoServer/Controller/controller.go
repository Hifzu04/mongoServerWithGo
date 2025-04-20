// controller.go
package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	model "github.com/Hifzu04/myMongoServer/Model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "netflix"
	collectionName = "watchlist"
)

var collection *mongo.Collection

// init connects to MongoDB and sets up the collection handle
func init() {
	// Load .env (optional)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠️  Controller.init: .env not found, using env vars")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("❌ Controller.init: MONGODB_URI is not set")
	}

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatalf("❌ Controller.init: cannot connect to MongoDB: %v", err)
	}

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Controller.init: ping to MongoDB failed: %v", err)
	}

	log.Println("✅ Connected to MongoDB")
	collection = client.Database(dbName).Collection(collectionName)
}

// Below your init(), still in controller.go

// writeJSON is a helper to send JSON + status code
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// InsertOneMovie handles POST /api/movie
func InsertOneMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Netflix
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, movie)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "insert failed"})
		return
	}
	movie.ID = res.InsertedID.(primitive.ObjectID)
	writeJSON(w, http.StatusCreated, movie)
}

// GetAllMovies handles GET /api/movies
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "fetch failed"})
		return
	}
	defer cursor.Close(ctx)

	var movies []model.Netflix
	for cursor.Next(ctx) {
		var m model.Netflix
		if err := cursor.Decode(&m); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "decode failed"})
			return
		}
		movies = append(movies, m)
	}
	writeJSON(w, http.StatusOK, movies)
}

// UpdateOneMovie handles PUT /api/movie/{id}
func UpdateOneMovie(w http.ResponseWriter, r *http.Request) {
	idHex := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"watched": true}})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteOneMovie handles DELETE /api/movie/{id}
func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	idHex := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// DeleteAllMovies handles DELETE /api/movies
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete-all failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"deletedCount": res.DeletedCount})
}
