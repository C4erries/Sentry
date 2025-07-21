package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Produce model to Kafka",
	Long: `This command produces model to Kafka. 
	You can specify the event type, user ID, and the number of model to send. For example:
		To produce 10 login model for user ID 456 from Russia:
  produce --type login --user_id 456 --count 10 --country RU`,
	Run: func(cmd *cobra.Command, args []string) {
		runProduce()
	},
}

var (
	payloadType     string
	baseEvent       model.BaseEvent
	userId          int
	count           int
	loginData       model.LoginData
	transactionData model.TransactionData
)

func init() {
	rootCmd.AddCommand(produceCmd)

	produceCmd.Flags().StringVarP(&payloadType, "type", "t", "", "event type (login, click, transaction)")

	produceCmd.Flags().IntVarP(&userId, "user_id", "u", 0, "user ID")
	produceCmd.Flags().IntVarP(&count, "count", "c", 1, "number of events to send")
	produceCmd.Flags().StringVar(&baseEvent.IP, "ip", "", "IP adress")
	produceCmd.Flags().StringVar(&baseEvent.GeoCountry, "country", "", "Geo Country")

	//Login
	produceCmd.Flags().StringVar(&loginData.Method, "method", "", "login method")
	produceCmd.Flags().BoolVar(&loginData.Success, "success", false, "login success")

	//Transaction
	produceCmd.Flags().Float64Var(&transactionData.Amount, "amount", 0, "transaction amount")
	produceCmd.Flags().StringVar(&transactionData.Currency, "currency", "", "transaction currency")

	produceCmd.MarkFlagRequired("type")
	produceCmd.MarkFlagRequired("user_id")
}

func runProduce() {
	baseEvent.EventType = model.EventType(payloadType)
	baseEvent.UserID = "#" + strconv.Itoa(userId)

	p, err := kafka.NewProducer([]string{os.Getenv("KAFKA_ADDR")}, "events_topic")
	if err != nil {
		slog.Error("create new producer error", slog.Any("err", err))
		os.Exit(1)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	es := make([]*model.Event, 0, count)
	for i := 0; i < count; i++ {

		var data interface{}
		switch baseEvent.EventType {
		case model.EventLogin:
			data = loginData
		case model.EventTransaction:
			data = transactionData
		default:
			slog.ErrorContext(ctx, "data is not assignable.", slog.Any("err", fmt.Errorf("unknown event type: %s", baseEvent.EventType.String())))
			os.Exit(1)
		}

		currentEvent := baseEvent
		currentEvent.Timestamp = time.Now().UTC()

		e, err := model.NewEvent(currentEvent, data)
		if err != nil {
			slog.ErrorContext(ctx, "NewEvent creation failed", slog.Any("err", err))
			os.Exit(1)
			return
		}
		es = append(es, e)

	}

	if err = p.ProduceBatch(ctx, es...); err != nil {
		slog.ErrorContext(ctx, "produce batch failed", slog.Any("err", err))
		os.Exit(1)
	}
}
