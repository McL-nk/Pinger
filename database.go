package main

import (
	"context"
	"fmt"

	"log"

	"github.com/iverly/go-mcping/api/types"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

func getServers(client *mongo.Client) (results []Serverstruct) {

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

func updateServer(client *mongo.Client, server Serverstruct, result *types.PingResponse) {

	coll := client.Database("bamb").Collection("servers")
	var id = server.ID

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{}

	if result.Protocol == -110 {

		update = bson.D{{Key: "$set", Value: bson.D{{Key: "online_players", Value: 0}, {Key: "max_players", Value: 0}, {Key: "online", Value: false}}}}

	} else {

		update = bson.D{{Key: "$set", Value: bson.D{{Key: "online_players", Value: result.PlayerCount.Online}, {Key: "max_players", Value: result.PlayerCount.Max}, {Key: "online", Value: true}}}}

	}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("Update server for " + server.Ip)

}
