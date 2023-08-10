package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/McL-nk/Pinger/database"
	"github.com/McL-nk/Pinger/services"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Check the .env stuff
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGO_URI"); uri == "" {
		log.Fatal("No MONGO_URI env var set")
	}

	// Create the mongo client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	go services.UpdateChannels(client)

	// Create cron scheduler
	s := gocron.NewScheduler(time.UTC)

	fmt.Println("Ready!")

	//	Setup the cron job to ping the servers every 5 min
	s.Every(5).Minutes().Do(func() {
		results := database.GetServers(client)
		services.BeginPing(results, client)
	})

	//	Start the cron jobs
	s.StartBlocking()

}
