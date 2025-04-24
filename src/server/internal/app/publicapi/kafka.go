package publicapi

import (
	kafkahandler "clusterlizer/internal/handler/kafka"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func startKafkaConsumers(cfg KafkaConsumer, log *zap.SugaredLogger) error {

	documentSaverConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.DocumentSaver.URL},
		Topic:     cfg.DocumentSaver.Topic,
		GroupID:   cfg.DocumentSaver.Group,
		Partition: cfg.DocumentSaver.Partition,
		MaxBytes:  cfg.DocumentSaver.MaxBytes, // 10MB
	})

	handler := kafkahandler.New(documentSaverConsumer, kafkahandler.Config{
		Topic:    cfg.DocumentSaver.Topic,
		URL:      cfg.DocumentSaver.URL,
		MaxBytes: cfg.DocumentSaver.MaxBytes},
		log,
	)
	if err := handler.DocumentSaver(); err != nil {
		return fmt.Errorf("kafka documer saver: %w", err)
	}
	return nil

}
