package config

import (
	"strings"

	"github.com/c4erries/Sentry/internal/logger"
	"github.com/c4erries/Sentry/internal/redis"
	"github.com/c4erries/Sentry/internal/storage"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}
	Kafka struct {
		Brokers         []string `mapstructure:"brokers"`
		EventTopic      string   `mapstructure:"event_topic"`
		ConsumerGroupID string   `mapstructure:"consumer_group_id"`
	}
	Redis    redis.Config
	Postgres storage.Config
	Log      logger.Config
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	// 1. Defaults
	setDefaults(v)

	// 2. ENV
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3. Файл конфигурации (опционально)
	// Ищем config.{yaml,json,toml}
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/sentry/")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// 4. Flags (Cobra)
	pflag.Int("server.port", v.GetInt("server.port"), "HTTP server port")
	pflag.String("server.host", v.GetString("server.host"), "HTTP server host")
	v.BindPFlags(pflag.CommandLine)

	// 5. Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
