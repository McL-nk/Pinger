package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Serverstruct struct {
	ID             primitive.ObjectID `bson:"_id"`
	Ip             string
	Version        string
	Online         bool
	Uptime         []interface{}
	Player_numbers []interface{}
	Max_players    int
	Online_players int
}
