package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Serverstruct struct {
	ID             primitive.ObjectID `bson:"_id"`
	Ip             string
	Version        string
	Uptime         []interface{}
	Player_numbers []interface{}
}
