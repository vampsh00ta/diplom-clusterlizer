package document

import (
	"clusterlizer/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ServiceImpl struct {
	producer *kafka.Writer
	storage  storage.Storage
	log      *zap.SugaredLogger
}

type Service interface {
	SendToBroker(ctx context.Context, data any) error
}

func New(
	producer *kafka.Writer,
	storage storage.Storage,
	log *zap.SugaredLogger) *ServiceImpl {
	return &ServiceImpl{
		producer: producer,
		storage:  storage,
		log:      log,
	}
}

func (s *ServiceImpl) SendToBroker(ctx context.Context, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("send document marshal: %w", err)
	}
	if err = s.producer.WriteMessages(ctx, kafka.Message{Value: dataBytes}); err != nil {
		return fmt.Errorf("send producer data: %w", err)
	}
	return nil
}
