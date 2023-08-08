package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Minecraft server model

type Serverstruct struct {
	ID             primitive.ObjectID `bson:"_id"`
	Ip             string
	Version        string
	Online         bool
	Uptime         []interface{}
	Player_numbers []PlayerNumbers
	Max_players    int
	Online_players int
}

type PlayerNumbers struct {
	Online int
	Time   int
}

// Guild model

type Guildstruct struct {
	ID      string `bson:"_id"`
	Server  server
	Name    string
	Premium int
}

type server struct {
	Id              string
	Status_channels status_channels
}

type status_channels struct {
	Players_online int
	Server_online  int
}
