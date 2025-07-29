package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/Sentry/internal/config"
	"github.com/c4erries/Sentry/internal/dispatcher"
	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/processor"
	"github.com/c4erries/Sentry/internal/redis"
	"github.com/c4erries/Sentry/internal/storage"
	"github.com/c4erries/Sentry/internal/worker"
	go_redis "github.com/redis/go-redis/v9"
)

func Serve(cfg *config.Config) {
	slog.Info("Config", slog.Any("config", cfg))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := go_redis.NewClient(&go_redis.Options{
		Addr: cfg.Redis.Addr,
	})
	wrappedRedis := redis.NewAdapter(redisClient)

	postgres, err := storage.NewStorage(&cfg.Postgres)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to db", slog.Any("err", err))
		os.Exit(1)
	}

	storageSink := dispatcher.NewStorageSink(postgres.Events, postgres.Alerts)

	alertDispatcher := dispatcher.NewAlertDispatcher(
		[]dispatcher.AlertSink{
			dispatcher.NewCmdSink(),
			dispatcher.NewRedisSink(wrappedRedis, redis.NewRedisPubSub(redisClient, "alerts")),
			storageSink,
		},
		10,
	)
	go alertDispatcher.Run(ctx)

	eventDispatcher := dispatcher.NewEventDispatcher(
		[]dispatcher.EventSink{
			storageSink,
		},
		10,
	)
	go eventDispatcher.Run(ctx)

	processor, err := processor.NewProcessor(wrappedRedis, eventDispatcher, alertDispatcher)
	if err != nil {
		slog.ErrorContext(ctx, "cannot create processor", slog.Any("err", err))
		os.Exit(1)
	}

	jobs := make(chan *kafka.KafkaEvent, 100)
	wg := worker.StartPool(ctx, jobs, processor, 5)

	kr, err := kafka.NewReader(
		cfg.Kafka.Brokers,
		cfg.Kafka.EventTopic,
		cfg.Kafka.ConsumerGroupID,
	)
	if err != nil {
		slog.ErrorContext(ctx, "cannot create kafka reader", slog.Any("err", err))
		os.Exit(1)
	}
	consumer, err := kafka.NewConsumer(kr)
	if err != nil {
		slog.ErrorContext(ctx, "can`t create kafka consumer", slog.Any("err", err))
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
	if err := postgres.Close(); err != nil {
		slog.Error("postgres close error", slog.Any("err", err))
	}
	wg.Wait()
	slog.InfoContext(ctx, "Service stoped.")
}
