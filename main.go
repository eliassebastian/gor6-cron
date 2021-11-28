package main

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eliassebastian/gor6-cron/internal/pubsub"
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

	producer, err := pubsub.NewKafkaConnection(context.Background(), "ubisoft-topic")
	if err != nil {
		return nil, err
	}

	client := ubisoft.CreateConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	srv := &Server{
		kafka:     producer,
		scheduler: scheduler,
		ubisoft:   client,
		doneC:     make(chan struct{}),
		closeC:    make(chan struct{}),
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
			producer.Producer.Close()
			scheduler.Stop()
			client.Stop()

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
	kafka     *pubsub.Producer
	scheduler *gocron.Scheduler
	ubisoft   *ubisoft.UbisoftConfig
	doneC     chan struct{}
	closeC    chan struct{}
}

func (s *Server) ListenAndServe() error {
	log.Println(":::::: ListenAndServer Func")

	err := s.ubisoft.Connect(context.Background(), s.kafka)
	if err != nil {
		return errors.New("ubisoft connect failure")
	}

	//TODO initiate cron job every 2hr45min
	//s.scheduler.Every("2h45m").Do()
	//s.scheduler.Every("1m").Do()
	s.scheduler.StartBlocking()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println(":::::: Shutting Down Server")

	close(s.closeC)
	for {
		select {
		case <-ctx.Done():
			return errors.New("Context.Done Error")
		case <-s.doneC:
			return nil
		}
	}
}
