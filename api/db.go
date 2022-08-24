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
	log.Printf("mongourl:%s", mongoUrl)
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		// ApplyURI(mongoUrl).
		ApplyURI("mongodb://localhost:27017").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("connect error:%s", err)
	}
	coll := client.Database("sample_training").Collection("zips")
	res := coll.FindOne(context.Background(), bson.M{
		"_id": bson.ObjectIdHex("5c8eccc1caa187d17ca76050"),
	})
	if res.Err() != nil {
		log.Fatalf("find one error:%s", res.Err())
	}
	doc := struct {
		City  string `bson:"city"`
		Pop   int32  `bson:"pop"`
		State string `bson:"state"`
		Zip   string `bson:"zip"`
	}{}
	/*
		err := coll.Find(bson.M{
			"_id": bson.ObjectIdHex("5c8eccc1caa187d17ca76050"),
		}).One(&doc)
	*/
	err = res.Decode(&doc)
	if err != nil {
		log.Fatalf("decode error:%s", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "fond one, city:%s, pop:%d, state:%s, zip:%s", doc.City, doc.Pop, doc.State, doc.Zip)
}
