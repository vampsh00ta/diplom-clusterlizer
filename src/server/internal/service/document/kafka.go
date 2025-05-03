package document

import (
	"context"
	"encoding/json"
	"fmt"

	"clusterlizer/internal/storage"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewKafka(
	producer *kafka.Writer,
	storage storage.Storage,
	log *zap.SugaredLogger,
) *KafkaImpl {
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

func (s *KafkaImpl) SendDocumentNames(ctx context.Context, params SendDocumentParams) error {
	s.log.Info("send to broker")

	dataBytes, err := json.Marshal(params.Keys)
	if err != nil {
		s.log.Errorf("send document marshal: %w", err)

		return fmt.Errorf("send document marshal: %w", err)
	}
	if err := s.producer.WriteMessages(ctx, kafka.Message{Value: dataBytes}); err != nil {
		s.log.Errorf("send producer document names: %w", err)

		return fmt.Errorf("send producer document names: %w", err)
	}
	return nil
}
