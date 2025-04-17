package document

import (
	"clusterlizer/internal/storage"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewKafka(
	producer *kafka.Writer,
	storage storage.Storage,
	log *zap.SugaredLogger) *KafkaImpl {
	return &KafkaImpl{
		producer: producer,
		storage:  storage,
		log:      log,
	}
}

type KafkaImpl struct {
	producer *kafka.Writer
	storage  storage.Storage
	log      *zap.SugaredLogger
}

func (s *KafkaImpl) SendToBroker(ctx context.Context, data any) error {
	s.log.Info("send to broker")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		s.log.Errorf("send document marshal: %w", err)

		return fmt.Errorf("send document marshal: %w", err)
	}
	if err := s.producer.WriteMessages(ctx, kafka.Message{Value: dataBytes}); err != nil {
		s.log.Errorf("send producer data: %w", err)

		return fmt.Errorf("send producer data: %w", err)
	}
	return nil
}
