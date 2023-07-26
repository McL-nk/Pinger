package main

import (
	"github.com/iverly/go-mcping/api/types"
	"github.com/iverly/go-mcping/mcping"
	"go.mongodb.org/mongo-driver/mongo"

	"sync"
)

func ping(pinger types.Pinger, server Serverstruct, wg *sync.WaitGroup, client *mongo.Client) {

	defer wg.Done()

	result, err := pinger.Ping(server.Ip, 25565)

	if err != nil {

		var fallback = types.PingResponse{Latency: 0, // Latency between you and the server
			PlayerCount: types.PlayerCount{Max: 0, Online: 0}, // Players count information of the server
			Protocol:    -110,                                 // Protocol number of the server
			Favicon:     "null",                               // Favicon in base64 of the server
			Motd:        "null",                               // Motd of the server without color
			Version:     "null",                               // Version of the server
			Sample:      nil}

		updateServer(client, server, &fallback)
	}

	updateServer(client, server, result)
}

func beginPing(servers []Serverstruct, client *mongo.Client) {
	var wg sync.WaitGroup

	pinger := mcping.NewPinger()

	wg.Add(len(servers))

	for _, server := range servers {

		go ping(pinger, server, &wg, client)
	}

	wg.Wait()
}
