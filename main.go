package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/eliassebastian/gor6-cron/internal/cache"
)

func main() {
	log.Println("::::: R6 CRON STARTING")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//TODO Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())

	//connect to redis
	err := cache.InitCache(ctx)
	if err != nil {
		log.Println(err)
	}
	//create initial cron job

	//set up scheduler
	osCall := <-c
	log.Printf("::::: R6 CRON SIGNAL RECEIVED: %+v", osCall)
	cancel()
}
