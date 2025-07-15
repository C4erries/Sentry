package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
	}
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}
	Postgres struct {
		DSN string `mapstructure:"dsn"`
	}
	Log struct {
		Level string `mapstructure:"level"`
	}
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	// 1. Defaults
	setDefaults(v)

	// 2. ENV
	v.SetEnvPrefix("SENTRY")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3. Файл конфигурации (опционально)
	// Ищем config.{yaml,json,toml} в cwd и /etc/sentry/
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
	// …другие флаги…
	//pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// 5. Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
