package services

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/McL-nk/Pinger/database"
	"github.com/McL-nk/Pinger/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct {
	Item string
	Done bool
}

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

type changeEvent struct {
	ID            changeID            `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  models.Serverstruct `bson:"fullDocument"`
	DocumentKey   documentKey         `bson:"documentKey"`
	Ns            namespace           `bson:"ns"`
}

func UpdateChannels(client *mongo.Client) {
	coll := client.Database("bamb").Collection("servers")
	// open a change stream with an empty pipeline parameter

	pipeline := mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "operationType", Value: "update"}}}}}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	changeStream, err := coll.Watch(context.TODO(), pipeline, opts)
	if err != nil {
		panic(err)
	}
	defer changeStream.Close(context.TODO())
	// iterate over the cursor to print the change-stream events
	for changeStream.Next(context.TODO()) {

		var changeEvent changeEvent
		changeStream.Decode(&changeEvent)

		key := changeEvent.DocumentKey
		updateChannel(client, fmt.Sprintf(key.ID.Hex()), changeEvent)
	}

}

func updateChannel(client *mongo.Client, key string, document changeEvent) {

	guilds := database.GetGuilds(client, key)
	httpclient := &http.Client{}

	for _, e := range guilds {
		if e.Server.Status_channels.Server_online != 0 {
			body := []byte(`{}`)
			if document.FullDocument.Online == true {
				body = []byte(`{
					"name": "🟢 Server Online"
				}`)
			} else {
				body = []byte(`{
					"name": "🔴 Server offline"
				}`)
			}

			patchurl := fmt.Sprintf("https://discord.com/api/v10/channels/%d", e.Server.Status_channels.Server_online)
			r, err := http.NewRequest("PATCH", patchurl, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			r.Header.Add("authorization", fmt.Sprintf("Bot %v", os.Getenv("BOT_TOKEN")))
			r.Header.Add("accept", "/")
			r.Header.Add("authority", "discord.com")
			r.Header.Add("content-type", "application/json")

			res, err := httpclient.Do(r)

			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

		}
		if e.Server.Status_channels.Players_online != 0 {
			body := []byte(`{}`)
			if document.FullDocument.Online == true {

				bod := fmt.Sprintf(`{"name": "👥 Players Online: %v"}`, document.FullDocument.Online_players)
				body = []byte(bod)

			} else {

				body = []byte(`{"name": "👥 Players Online: 0"}`)

			}

			patchurl := fmt.Sprintf("https://discord.com/api/v10/channels/%d", e.Server.Status_channels.Players_online)
			r, err := http.NewRequest("PATCH", patchurl, bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
			r.Header.Add("authorization", fmt.Sprintf("Bot %v", os.Getenv("BOT_TOKEN")))
			r.Header.Add("accept", "/")
			r.Header.Add("authority", "discord.com")
			r.Header.Add("content-type", "application/json")

			res, err := httpclient.Do(r)

			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

		}
	}

}
