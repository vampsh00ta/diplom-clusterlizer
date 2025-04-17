package publicapi

import (
	documentsrvc "clusterlizer/internal/service/document"

	"github.com/segmentio/kafka-go"

	psqlrep "clusterlizer/internal/storage/postgres"
	"clusterlizer/pkg/pgxclient"
	"context"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	sugar := logger.Sugar()
	return sugar
}
func Run(cfg *Config) {
	ctx := context.Background()
	defer ctx.Done()
	// Logger
	logger := NewLogger()
	logger.Info("starting pg...")

	// PostgresQL
	pg, err := pgxclient.New(ctx, 5,
		pgxclient.Config{
			Name:     cfg.PG.Name,
			Username: cfg.PG.Username,
			Port:     cfg.PG.Port,
			Password: cfg.PG.Password,
			Host:     cfg.PG.Host,
		})

	// Kafka Producer
	documentSenderProducer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Producer.DocumentSender.URL),
		Topic:                  cfg.Kafka.Producer.DocumentSender.Topic,
		RequiredAcks:           kafka.RequireOne,
		Balancer:               &kafka.LeastBytes{},
		Async:                  true,
		AllowAutoTopicCreation: true,
	}

	if err != nil {
		logger.Fatal("user - Run - postgres.New: %v", zap.Error(err))
	}
	defer pg.Close()

	//Storage

	logger.Info("starting storage...")

	storage := psqlrep.New(pg)

	logger.Info("starting services...")

	//Services

	documentImpl := documentsrvc.NewKafka(
		documentSenderProducer,
		storage,
		logger,
	)

	// HTTP server
	logger.Info("starting HTTP server...")

	app := registerHTPP(cfg, logger, documentImpl)

	// Kafka consumers
	logger.Info("starting kafka consumer...")

	// if err = startKafkaConsumers(cfg.Kafka.Consumer, logger); err != nil {
	// 	logger.Fatal("kafka consumer: %w", zap.Error(err))
	// }

	if err = app.Listen(":" + cfg.App.Port); err != nil {
		logger.Fatal("Ошибка запуска сервера: %v", zap.Error(err))
	}
}
