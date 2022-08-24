package handler

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDb(t *testing.T) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		// ApplyURI(mongoUrl).
		// ApplyURI("mongodb+srv://kaiya:kaiyaxiong123@cluster0.kowylj0.mongodb.net/?retryWrites=true&w=majority").
		ApplyURI("mongodb://localhost:27017").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("connect error:%s", err)
	}
	databases, _ := client.ListDatabaseNames(ctx, bson.M{})
	log.Printf("%v", databases)
	collections, _ := client.Database("sample_training").ListCollectionNames(ctx, bson.M{})
	log.Printf("%v", collections)
	coll := client.Database("sample_training").Collection("zips")
	res := coll.FindOne(context.Background(), bson.M{
		"_id": bson.ObjectIdHex("5c8eccc1caa187d17ca76050"),
	})
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

}
