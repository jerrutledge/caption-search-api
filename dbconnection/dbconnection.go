package dbconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jerrutledge/caption-search-api/episode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Collection, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	collection := client.Database("caption-search").Collection("episodes")
	return collection, err
}

func HelloResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func SearchResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := "Wizard"
	coll, err := Connect()
	if err != nil {
		ReturnError(w)
		return
	}
	var response = episode.SearchResults{Err: false}
	err, response.Results = episode.Search(coll, query)
	if err != nil {
		ReturnError(w)
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		ReturnError(w)
		return
	}
	w.Write(data)
	return
}

func ReturnError(w http.ResponseWriter) {
	var response = episode.SearchResults{Err: true}
	data, err := json.Marshal(response)
	if err != nil {
		// TODO: handle error
	}
	w.Write(data)
	return

}