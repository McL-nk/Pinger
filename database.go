package main

import (
	"context"

	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getServers() (results []Serverstruct) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGO_URI"); uri == "" {
		log.Fatal("No MONGO_URI env var set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("bamb").Collection("servers")

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result Serverstruct
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}
