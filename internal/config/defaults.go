package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)

	v.SetDefault("kafka.brokers", []string{"kafka:9092"})
	v.SetDefault("kafka.event_topic", "events_topic")
	v.SetDefault("kafka.consumer_group_id", "sentry-core")

	v.SetDefault("redis.addr", "kafka:6379")
	v.SetDefault("redis.db", 0)

	v.SetDefault("postgres.dsn", "postgres://sentry:sentrypass@postgres:5432/sentry_db?sslmode=disable")

	v.SetDefault("log.level", "info")
}
