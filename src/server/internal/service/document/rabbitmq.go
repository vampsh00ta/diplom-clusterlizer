package document

import (
	"context"
	"encoding/json"
	"fmt"

	"clusterlizer/internal/storage"

	rabbitmq "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQConfig struct {
	Exchange string
	Key      string
}

func NewRabbiqMQ(
	cfg RabbitMQConfig,
	producer *rabbitmq.Channel,
	storage storage.Storage,
	log *zap.SugaredLogger,
) *RabbitMQImpl {
	return &RabbitMQImpl{
		cfg:      cfg,
		producer: producer,
		storage:  storage,
		log:      log,
	}
}

type RabbitMQImpl struct {
	cfg      RabbitMQConfig
	producer *rabbitmq.Channel
	storage  storage.Storage
	log      *zap.SugaredLogger
}

func (s *RabbitMQImpl) SendDocumentNames(ctx context.Context, params SendDocumentParams) error {
	s.log.Info("send to rabbitmq")

	dataBytes, err := json.Marshal(params)
	if err != nil {
		s.log.Errorf("send document marshal: %w", err)

		return fmt.Errorf("send document marshal: %w", err)
	}
	if err := s.producer.PublishWithContext(ctx,
		s.cfg.Exchange,
		s.cfg.Key,
		false,
		false,
		rabbitmq.Publishing{
			Body: dataBytes,
		}); err != nil {
		s.log.Errorf("send producer document names: %w", err)

		return fmt.Errorf("send producer document names: %w", err)
	}
	return nil
}
