package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/c4erries/Sentry/internal/kafka/mocks"
	"github.com/c4erries/Sentry/internal/testutils"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConsumer_Start(t *testing.T) {
	validEventRaw := testutils.ValidEventRaw(t)
	invalidEventRaw := testutils.InvalidEventRaw(t)

	type args struct {
		ctx  context.Context
		msgs []kafka.Message // входящие Kafka-сообщения
	}
	tests := []struct {
		name       string
		args       args
		wantEvents int // сколько сообщений ожидаем получить в out
	}{
		{
			name: "valid event",
			args: args{
				ctx: context.Background(),
				msgs: []kafka.Message{
					{
						Value: validEventRaw,
					},
				},
			},
			wantEvents: 1,
		},
		{
			name: "invalid event",
			args: args{
				ctx: context.Background(),
				msgs: []kafka.Message{
					{
						Value: invalidEventRaw,
					},
				},
			},
			wantEvents: 0,
		},
		{
			name: "one valid, one invalid",
			args: args{
				ctx: context.Background(),
				msgs: []kafka.Message{
					{
						Value: invalidEventRaw,
					},
					{
						Value: validEventRaw,
					},
				},
			},
			wantEvents: 1,
		},
		{
			name: "invalid json",
			args: args{
				ctx: context.Background(),
				msgs: []kafka.Message{
					{
						Value: []byte(`invalid json`),
					},
				},
			},
			wantEvents: 0,
		},
		{
			name: "context cancelled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				msgs: []kafka.Message{},
			},
			wantEvents: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReader := new(mocks.Reader)

			// Настраиваем mockReader.ReadMessage чтобы возвращать сообщения по очереди из tt.args.msgs,
			// а потом — ошибку context.Canceled чтобы остановить цикл
			calls := 0
			mockReader.On("ReadMessage", mock.Anything).Return(func(ctx context.Context) (kafka.Message, error) {
				if calls < len(tt.args.msgs) {
					msg := tt.args.msgs[calls]
					calls++
					return msg, nil
				}
				return kafka.Message{}, context.Canceled
			}).Maybe()

			mockReader.On("CommitMessages", mock.Anything, mock.Anything).Return(nil).Maybe()
			mockReader.On("Close").Return(nil).Once()

			c := &Consumer{
				reader: mockReader,
			}

			out := make(chan *KafkaEvent, 10)
			done := make(chan struct{})

			go func() {
				c.Start(tt.args.ctx, out)
				close(done) // сигнализируем, что consumer завершился
			}()

			// Проверяем сколько сообщений пришло в out
			received := 0
			timeout := time.After(200 * time.Millisecond)

		loop:
			for {
				select {
				case <-out:
					received++
				case <-timeout:
					break loop
				}
			}

			<-done // ждём завершения consumer, чтобы Close() точно вызвался

			assert.Equal(t, tt.wantEvents, received)
			mockReader.AssertExpectations(t)
		})
	}
}
