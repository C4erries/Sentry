package cmd

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/c4erries/Sentry/internal/kafka"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Produce model to Kafka",
	Long: `This command produces model to Kafka. 
	You can specify the event type, user ID, and the number of model to send. For example:
		To produce 10 login model for user ID 456:
  produce --type login --user_id 456 --count 10`,
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
	baseEvent.UserId = "#" + strconv.Itoa(userId)

	p, err := kafka.NewProducer([]string{"0.0.0.0:9092"}, "events_topic")
	if err != nil {
		log.Fatalf("create new producer error: %v", err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	es := make([]model.Event, 0, count)
	for i := 0; i < count; i++ {

		var data interface{}
		switch baseEvent.EventType {
		case model.EventLogin:
			data = loginData
		case model.EventTransaction:
			data = transactionData
		default:
			log.Fatalf("data is not assignable for that event: %v", baseEvent.EventType)
		}

		currentEvent := baseEvent
		currentEvent.Timestamp = time.Now().UTC()

		e := model.Event{
			BaseEvent: currentEvent,
			ID:        uuid.NewString(),
			Data:      data,
		}

		es = append(es, e)

	}

	if err = p.ProduceBatch(ctx, es...); err != nil {
		log.Fatalf("produce error: %v", err)
	}
}
