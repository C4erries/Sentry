package main

import (
	"github.com/c4erries/Sentry/internal/app"
	"github.com/spf13/cobra"
)

// serveCmd представляет команду запуска основного приложения
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск основного сервиса Sentry",
	Long: `Команда запускает основной сервис:
- HTTP сервер
- Kafka consumer (с обработкой событий)
- Инициализацию Redis, PostgreSQL и др.

Пример использования:
  sentry serve
`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve() {
	app.Serve(cfg)
}
