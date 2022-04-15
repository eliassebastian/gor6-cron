package main

import (
	"context"
	"errors"
	"github.com/eliassebastian/gor6-cron/internal/rabbitmq"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eliassebastian/gor6-cron/internal/ubisoft"
)

func main() {
	log.Println(":::::: Running MAIN")

	errC, err := run()
	if err != nil {
		log.Fatalf("Error running %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("error running %s", err)
	}
}

func run() (<-chan error, error) {
	log.Println(":::::: R6 CRON STARTING")

	//producer := kafka.NewKafkaWriter("ubisoft-topic")
	producer, err := rabbitmq.NewConnection()
	if err != nil {
		return nil, err
	}

	client := ubisoft.CreateConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	srv := &Server{
		//kafka:     producer,
		rabbitmq:  producer,
		scheduler: scheduler,
		ubisoft:   client,
		doneC:     make(chan struct{}),
	}

	errC := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		log.Println(":::::: Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			err := producer.Close()
			if err != nil {
				log.Println("Failed to close Kafka Connection")
			}
			client.Stop()
			scheduler.Stop()

			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Println(":::::: Shutdown Finished")
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type Server struct {
	//kafka     *kafka.Producer
	rabbitmq  *rabbitmq.RabbitConfig
	scheduler *gocron.Scheduler
	ubisoft   *ubisoft.UbisoftConfig
	doneC     chan struct{}
}

func (s *Server) ListenAndServe() error {
	log.Println(":::::: ListenAndServer")
	//TODO initiate cron job every 2hr45min
	//s.scheduler.Every("2h45m").Do()
	job, err := s.scheduler.Every("10m").Do(func(con *rabbitmq.RabbitConfig) {
		err := s.ubisoft.Connect(context.Background(), con)
		if err != nil {
			log.Println("Job Error", err)
		}
		log.Println("SUCCESS UBI")
	}, s.rabbitmq)
	log.Println(job, err)
	s.scheduler.StartBlocking()
	log.Println("Scheduler Stopped")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println(":::::: Shutting Down Server")
	for {
		select {
		case <-ctx.Done():
			return errors.New("Context.Done Error")
		case <-s.doneC:
			return nil
		}
	}
}
