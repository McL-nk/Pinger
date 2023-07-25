package main

import (
	"fmt"

	"github.com/iverly/go-mcping/api/types"
	"github.com/iverly/go-mcping/mcping"

	"sync"
)

func ping(pinger types.Pinger, ip string, wg *sync.WaitGroup) {

	defer wg.Done()

	response, err := pinger.Ping(ip, 25565)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)

}

func beginPing(servers []Serverstruct) {
	var wg sync.WaitGroup

	pinger := mcping.NewPinger()

	wg.Add(len(servers))

	for _, server := range servers {

		go ping(pinger, server.Ip, &wg)
	}

	wg.Wait()
}
