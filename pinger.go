package main

import (
	"fmt"

	"github.com/iverly/go-mcping/api/types"
	"github.com/iverly/go-mcping/mcping"
	"go.mongodb.org/mongo-driver/mongo"

	"sync"
)

func ping(pinger types.Pinger, server Serverstruct, wg *sync.WaitGroup, client *mongo.Client) {

	defer wg.Done()

	result, err := pinger.Ping(server.Ip, 25565)

	if err != nil {
		fmt.Printf("server is offline")
		var fallback = types.PingResponse{Latency: 0,
			PlayerCount: types.PlayerCount{Max: 0, Online: 0},
			Protocol:    -110, // -110 protocol indicates that the ping was failed and this is the fallback response
			Favicon:     "null",
			Motd:        "null",
			Version:     "null",
			Sample:      nil}

		updateServer(client, server, &fallback)
		return
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
