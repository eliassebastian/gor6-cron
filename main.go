package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("::::: R6 CRON STARTING")

	//connect to Kafka

	//connect to redis

	//create initial cron job

	//set up scheduler

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	osCall := <-c
	log.Printf("::::: R6 CRON SIGNAL RECEIVED: %+v", osCall)
}
