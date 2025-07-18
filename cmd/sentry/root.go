package main

import (
	"log"
	"os"

	"github.com/c4erries/Sentry/internal/config"
	"github.com/c4erries/Sentry/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sentry",
	Short: "Sentry anomaly detection service",
	Long: `Sentry — это инструмент для детектирования и обработки
пользовательских событий в режиме реального времени.

Доступные команды:
  serve     Запускает HTTP‑сервер, Kafka consumer, Redis и т.д.
  migrate   Выполняет миграции базы данных
  produce   Генерирует и отправляет тестовые события в Kafka

Для запуска сервиса просто выполните:
  sentry
или
  sentry serve
`,
	// Если не указана подкоманда, по умолчанию выполняем serve
	RunE: func(cmd *cobra.Command, args []string) error {
		// явно вызываем логику serveCmd
		return serveCmd.RunE(serveCmd, args)
	},
}

func Execute() error {
	// Если нет аргументов, подставляем "serve"
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "serve")
	}
	return rootCmd.Execute()
}

var cfg *config.Config

func init() {
	rootCmd.PersistentFlags().AddFlagSet(pflag.CommandLine)

	cobra.OnInitialize(func() {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			log.Fatalf("cannot load config: %v", err)
		}
		if err := logger.Init(&cfg.Log); err != nil {
			log.Fatalf("cannot initialize logger: %v", err)
		}
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Sentry.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
