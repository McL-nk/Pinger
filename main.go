package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minutes().Do(func() {
		results := getServers()
		beginPing(results)
	})

	s.StartBlocking()
}
