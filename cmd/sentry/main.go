package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/Sentry/internal/events"
	"github.com/c4erries/Sentry/internal/handler"
	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan events.Event, 100)

	processor, err := handler.NewProcessor()
	if err != nil {
		log.Fatalf("cannot create processor: %v", err)
	}
	wg := worker.StartPool(ctx, jobs, processor, 5)

	consumer, err := kafka.NewConsumer([]string{"kafka:29092"}, "events_topic", "sentry-core")
	if err != nil {
		log.Fatalf("cannot create consumer: %v", err)
	}
	go func() {
		consumer.Start(ctx, jobs)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	cancel()
	close(jobs)
	wg.Wait()
	log.Println("Service stopped.")
}
