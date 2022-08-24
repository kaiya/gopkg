package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandlerDb(w http.ResponseWriter, r *http.Request) {
	mongoUrl := os.Getenv("MONGODB_URI")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoUrl).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("connect error:%s", err)
	}
	zip := r.FormValue("zip")
	coll := client.Database("sample_training").Collection("zips")
	res := coll.FindOne(context.Background(), bson.M{
		"zip": zip,
	})
	if res.Err() != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "zip:%s not found", zip)
		return
	}
	doc := struct {
		City  string `bson:"city"`
		Pop   int32  `bson:"pop"`
		State string `bson:"state"`
		Zip   string `bson:"zip"`
	}{}
	err = res.Decode(&doc)
	if err != nil {
		log.Fatalf("decode error:%s", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "city:%s, pop:%d, state:%s, zip:%s", doc.City, doc.Pop, doc.State, doc.Zip)
}
