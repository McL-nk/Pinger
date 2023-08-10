package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/McL-nk/Pinger/database"
	"github.com/McL-nk/Pinger/models"
	"github.com/iverly/go-mcping/api/types"
	"github.com/iverly/go-mcping/mcping"
	"go.mongodb.org/mongo-driver/mongo"

	"sync"
)

func ping(server models.Serverstruct, wg *sync.WaitGroup, client *mongo.Client) {
	pinger := mcping.NewPinger()

	defer wg.Done()

	port := strings.Split(server.Ip, ":")
	ip := server.Ip
	port2 := int64(0)
	if len(port) > 1 {
		ip = port[0]
		port2, _ = strconv.ParseInt(port[1], 10, 0)
	} else {
		port2 = 0
	}

	portnum := 0
	if port2 < 65536 && port2 > 0 {
		portnum = int(port2)
	} else {
		portnum = 25565
	}
	fmt.Print("server ip for: ", server.Ip)
	fmt.Println(" ", uint16(portnum))
	result, err := pinger.Ping(ip, uint16(portnum))

	if err != nil {

		fmt.Printf("server is offline")
		var fallback = types.PingResponse{Latency: 0,
			PlayerCount: types.PlayerCount{Max: 0, Online: 0},
			Protocol:    -110, // -110 protocol indicates that the ping was failed and this is the fallback response
			Favicon:     "null",
			Motd:        "null",
			Version:     "null",
			Sample:      nil}

		database.UpdateServer(client, server, &fallback)
		return
	}

	database.UpdateServer(client, server, result)
}

func BeginPing(servers []models.Serverstruct, client *mongo.Client) {
	var wg sync.WaitGroup

	wg.Add(len(servers))

	for _, server := range servers {

		go ping(server, &wg, client)
	}

	wg.Wait()
}
