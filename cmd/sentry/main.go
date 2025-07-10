package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/Sentry/internal/alert"
	"github.com/c4erries/Sentry/internal/handler"
	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/redis"
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

	dispatcher := alert.NewDispatcher(
		[]alert.AlertSink{
			alert.NewCmdSink(),
			alert.NewRedisSink(wrappedRedis, redis.NewRedisPubSub(redisClient, "alerts")),
		},
		5,
	)

	go dispatcher.Run(ctx)

	processor, err := handler.NewProcessor(wrappedRedis, dispatcher)
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
	close(dispatcher.Chan)
	wg.Wait()
	log.Println("Service stopped.")
}
