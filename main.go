package main

import (
	"context"
	"github.com/eliassebastian/gor6-cron/internal/cache"
	"github.com/eliassebastian/gor6-cron/internal/ubisoft"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("::::: R6 CRON STARTING")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//TODO Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())

	//connect to redis
	conn, err := cache.InitCache(ctx)
	if conn != nil || err != nil {
		log.Println(err)
	}
	//create initial cron job
	client, _ := ubisoft.EstablishConn()

	//set up scheduler
	osCall := <-c
	log.Printf("::::: R6 CRON SIGNAL RECEIVED: %+v", osCall)
	cancel()
}
