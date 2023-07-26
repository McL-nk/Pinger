package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

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

	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minutes().Do(func() {
		results := getServers(client)
		beginPing(results, client)
	})

	s.StartBlocking()
}
