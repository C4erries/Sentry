package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/Sentry/internal/dispatcher"
	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/processor"
	"github.com/c4erries/Sentry/internal/redis"
	"github.com/c4erries/Sentry/internal/storage"
	"github.com/c4erries/Sentry/internal/worker"
	go_redis "github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := go_redis.NewClient(&go_redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	wrappedRedis := redis.NewAdapter(redisClient)

	postgres, err := storage.NewStorage(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	storageSink := dispatcher.NewStorageSink(postgres.Events, postgres.Alerts)

	alertDispatcher := dispatcher.NewAlertDispatcher(
		[]dispatcher.AlertSink{
			dispatcher.NewCmdSink(),
			dispatcher.NewRedisSink(wrappedRedis, redis.NewRedisPubSub(redisClient, "alerts")),
			storageSink,
		},
		5,
	)
	go alertDispatcher.Run(ctx)

	eventDispatcher := dispatcher.NewEventDispatcher(
		[]dispatcher.EventSink{
			storageSink,
		},
		5,
	)
	go eventDispatcher.Run(ctx)

	processor, err := processor.NewProcessor(wrappedRedis, eventDispatcher, alertDispatcher)
	if err != nil {
		log.Fatalf("cannot create processor: %v", err)
	}

	jobs := make(chan *kafka.KafkaEvent, 100)
	wg := worker.StartPool(ctx, jobs, processor, 5)

	consumer, err := kafka.NewConsumer([]string{os.Getenv("KAFKA_ADDR")}, "events_topic", "sentry-core")
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
	close(alertDispatcher.Chan)
	close(eventDispatcher.Chan)
	postgres.DB.Close()
	wg.Wait()
	log.Println("Service stopped.")
}
