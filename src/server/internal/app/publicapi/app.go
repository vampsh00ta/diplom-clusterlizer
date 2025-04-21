package publicapi

import (
	documentsrvc "clusterlizer/internal/service/document"
	requestsrvc "clusterlizer/internal/service/request"
	s3srvc "clusterlizer/internal/service/s3"
	psqlrep "clusterlizer/internal/storage/postgres"
	"clusterlizer/pkg/pgxclient"
	s3client "clusterlizer/pkg/s3"

	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/segmentio/kafka-go"

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
	////S3 client
	s3Session := s3.NewFromConfig(cfg.S3.Config)
	s3Client := s3client.NewClient(s3Session, cfg.S3.Bucket)

	// Kafka Producer
	documentSenderProducer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Producer.DocumentNameSender.URL),
		Topic:                  cfg.Kafka.Producer.DocumentNameSender.Topic,
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

	//Services
	logger.Info("starting services...")

	documentImpl := documentsrvc.NewKafka(
		documentSenderProducer,
		storage,
		logger,
	)
	requestImpl := requestsrvc.NewRequest(
		storage,
		logger,
	)
	s3Impl := s3srvc.New(
		logger,
		s3Client,
	)
	// HTTP server
	logger.Info("starting HTTP server...")

	app := registerHTPP(cfg, logger, documentImpl, requestImpl, s3Impl)

	// Kafka consumers
	logger.Info("starting kafka consumer...")

	// if err = startKafkaConsumers(cfg.Kafka.Consumer, logger); err != nil {
	// 	logger.Fatal("kafka consumer: %w", zap.Error(err))
	// }

	if err = app.Listen(":" + cfg.App.Port); err != nil {
		logger.Fatal("Ошибка запуска сервера: %v", zap.Error(err))
	}
}
