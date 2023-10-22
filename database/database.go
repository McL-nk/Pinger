package database

import (
	"context"
	"fmt"
	"time"

	"github.com/McL-nk/Pinger/models"

	"github.com/iverly/go-mcping/api/types"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetServers(client *mongo.Client) (results []models.Serverstruct) {

	coll := client.Database("bamb").Collection("servers")

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Print("Error: ", err)
	}

	for cursor.Next(context.TODO()) {
		var result models.Serverstruct
		if err := cursor.Decode(&result); err != nil {
			fmt.Print("Error: ", err)
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		fmt.Print("Error: ", err)
	}
	return results
}

func UpdateServer(client *mongo.Client, server models.Serverstruct, result *types.PingResponse) {

	coll := client.Database("bamb").Collection("servers")
	var id = server.ID

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{}

	if len(server.Player_numbers) >= 288 {

		server.Player_numbers = append(server.Player_numbers[:0], server.Player_numbers[1:]...)
	}

	if result.Protocol == -110 {

		server.Player_numbers = append(server.Player_numbers, models.PlayerNumbers{Online: 0, Time: int(time.Now().UnixMilli())})

		update = bson.D{{Key: "$set", Value: bson.D{{Key: "online_players", Value: 0}, {Key: "max_players", Value: 0}, {Key: "online", Value: false}, {Key: "player_numbers", Value: server.Player_numbers}}}}

	} else {

		server.Player_numbers = append(server.Player_numbers, models.PlayerNumbers{Online: result.PlayerCount.Online, Time: int(time.Now().UnixMilli())})

		update = bson.D{{Key: "$set", Value: bson.D{{Key: "online_players", Value: result.PlayerCount.Online}, {Key: "max_players", Value: result.PlayerCount.Max}, {Key: "online", Value: true}, {Key: "player_numbers", Value: server.Player_numbers}}}}

	}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Print("Error: ", err)
	}

	fmt.Println("Update server for " + server.Ip)

}

func GetGuilds(client *mongo.Client, key string) (results []models.Guildstruct) {

	filter := bson.D{{Key: "server.id", Value: key}}

	coll := client.Database("bamb").Collection("guilds")

	cursor, err := coll.Find(context.TODO(), filter)

	if err != nil {

		fmt.Print("Error: ", err)
	}

	for cursor.Next(context.TODO()) {
		var result models.Guildstruct

		if err := cursor.Decode(&result); err != nil {
			fmt.Print("Error: ", err)
		}

		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		fmt.Print("Error: ", err)
	}

	return results
}
