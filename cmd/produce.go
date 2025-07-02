package cmd

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/c4erries/Sentry/internal/events"
	kc "github.com/c4erries/Sentry/internal/kafkaclient"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Produce events to Kafka",
	Long: `This command produces events to Kafka. 
	You can specify the event type, user ID, and the number of events to send. For example:
		To produce 10 login events for user ID 456:
  produce --type login --user_id 456 --count 10`,
	Run: func(cmd *cobra.Command, args []string) {
		runProduce()
	},
}

var (
	eventType string
	userID    int
	count     int
)

func init() {
	rootCmd.AddCommand(produceCmd)

	produceCmd.Flags().StringVarP(&eventType, "type", "t", "", "event type (login, click, transaction)")
	produceCmd.Flags().IntVarP(&userID, "user_id", "u", 0, "user ID")
	produceCmd.Flags().IntVarP(&count, "count", "c", 1, "number of events to send")

	produceCmd.MarkFlagRequired("type")
	produceCmd.MarkFlagRequired("user_id")
}

func runProduce() {
	payloadType := events.EventType(eventType)
	if !payloadType.IsValid() {
		log.Fatalf("payload type is invalid")
	}
	p, err := kc.NewProducer([]string{"0.0.0.0:9092"}, "events_topic")
	if err != nil {
		log.Fatalf("create new producer error: %v", err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msgs := make([]kafka.Message, 0, count)
	for i := 0; i < count; i++ {

		payload := map[string]interface{}{
			"event_type": payloadType,
			"user_id":    userID,
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		}

		data, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("marshal error: %v", err)
		}

		msg := kafka.Message{
			Key:   []byte("user_" + strconv.Itoa(userID)),
			Value: data,
		}

		msgs = append(msgs, msg)

	}

	if err = p.Produce(ctx, msgs...); err != nil {
		log.Fatalf("produce error: %v", err)
	}
}
