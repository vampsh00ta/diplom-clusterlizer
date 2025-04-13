package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Handler struct {
	consumer *kafka.Reader
	log      *zap.SugaredLogger
}

type Config struct {
	URL      string ` yaml:"url"`
	Topic    string `yaml:"topic"`
	MaxBytes int    ` yaml:"max_bytes"`
}

func New(consumer *kafka.Reader, cfg Config, log *zap.SugaredLogger) Handler {
	return Handler{}
}
func (k Handler) DocumentSaver() error {
	ctx := context.Background()
	errorChan := make(chan error, 1)
	go func() {
		for err := range errorChan {
			k.log.Error("error occurred", zap.Error(err))
		}
	}()

	go func() {
		for {

			var err error
			msg, err := k.consumer.FetchMessage(ctx)
			if err != nil {
				errorChan <- err
				continue
			}
			//var data entity.Vacancy
			//if err = json.Unmarshal(msg.Value, &data); err != nil {
			//	errorChan <- err
			//}
			//
			//if err = k.job.AddSpecialiity(context.Background(), &data); err != nil {
			//	errorChan <- err
			//}
			//if err = k.consumer.CommitMessages(ctx, msg); err != nil {
			//	errorChan <- err
			//	continue
			//}
			k.log.Info("committed message",
				zap.Int("PartitionID", msg.Partition),
				zap.Int64("Offset", msg.Offset))
		}
	}()

	return nil
}
