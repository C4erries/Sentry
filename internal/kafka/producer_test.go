package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/c4erries/Sentry/internal/kafka/mocks"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/c4erries/Sentry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProducer_Produce(t *testing.T) {
	validEvent := testutils.ValidEvent(t)
	invalidEvent := testutils.InvalidEvent(t)

	type args struct {
		ctx context.Context
		e   *model.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx: context.Background(),
				e:   validEvent,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx: context.Background(),
				e:   invalidEvent,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			writer := new(mocks.Writer)
			if !tt.wantErr {
				writer.On("WriteMessages", mock.Anything, mock.Anything).Return(nil).Once()
			}

			producer, err := NewProducer(writer)
			require.NoError(t, err)

			err = producer.Produce(tt.args.ctx, tt.args.e)
			if tt.wantErr {
				assert.Error(t, err)
				writer.AssertNotCalled(t, "WriteMessages", mock.Anything, mock.Anything)
			} else {
				assert.NoError(t, err)
				writer.AssertCalled(t, "WriteMessages", mock.Anything, mock.Anything)
			}

			writer.AssertExpectations(t)
		})
	}
}

func TestProducer_ProduceBatch(t *testing.T) {
	validEvent := testutils.ValidEvent(t)
	invalidEvent := testutils.InvalidEvent(t)

	type args struct {
		ctx context.Context
		evs []*model.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid event",
			args{
				ctx: context.Background(),
				evs: []*model.Event{validEvent},
			},
			false,
		},
		{
			"invalid event",
			args{
				ctx: context.Background(),
				evs: []*model.Event{invalidEvent},
			},
			true,
		},
		{
			"valid&invalid event",
			args{
				ctx: context.Background(),
				evs: []*model.Event{validEvent, invalidEvent},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			writer := new(mocks.Writer)
			if !tt.wantErr {
				writer.On("WriteMessages", mock.Anything, mock.Anything).Return(nil).Once()
			}

			producer, err := NewProducer(writer)
			require.NoError(t, err)

			err = producer.ProduceBatch(tt.args.ctx, tt.args.evs...)
			if tt.wantErr {
				assert.Error(t, err)
				writer.AssertNotCalled(t, "WriteMessages", mock.Anything, mock.Anything)
			} else {
				assert.NoError(t, err)
				writer.AssertCalled(t, "WriteMessages", mock.Anything, mock.Anything)
			}

			writer.AssertExpectations(t)
		})
	}
}

func TestProducer_Close(t *testing.T) {
	tests := []struct {
		name       string
		closeError error
		wantErr    bool
	}{
		{
			name:       "writer.Close succeeds",
			closeError: nil,
			wantErr:    false,
		},
		{
			name:       "writer.Close fails",
			closeError: errors.New("close failed"),
			wantErr:    true,
		},
	}

	for _, ttOrig := range tests {
		tt := ttOrig
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			writer := new(mocks.Writer)
			writer.
				On("Close").
				Return(tt.closeError).
				Once()

			producer, err := NewProducer(writer)
			require.NoError(t, err)

			err = producer.Close()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			writer.AssertExpectations(t)
		})
	}
}
