package publicapi

import (
	"context"

	rabbitmqconsumer "clusterlizer/internal/handler/rabbitmq"
	documentsrvc "clusterlizer/internal/service/document"
	filesrvc "clusterlizer/internal/service/file"
	requestsrvc "clusterlizer/internal/service/request"

	s3srvc "clusterlizer/internal/service/s3"
	psqlrep "clusterlizer/internal/storage/postgres"
	"clusterlizer/pkg/pgxclient"
	s3client "clusterlizer/pkg/s3"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	rabbitmq "github.com/rabbitmq/amqp091-go"
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
	if err != nil {
		logger.Fatal("user - Run - postgres.New: %v", zap.Error(err))
	}
	defer pg.Close()
	////S3 client
	s3Session := s3.NewFromConfig(cfg.S3.Config)
	s3Client := s3client.NewClient(s3Session, cfg.S3.Bucket)

	//Kafka KafkaProducer
	//documentSenderProducer := &kafka.Writer{
	//	Addr:                   kafka.TCP(cfg.Kafka.Producer.DocumentNameSender.URL),
	//	Topic:                  cfg.Kafka.Producer.DocumentNameSender.Topic,
	//	RequiredAcks:           kafka.RequireOne,
	//	Balancer:               &kafka.LeastBytes{},
	//	Async:                  true,
	//	AllowAutoTopicCreation: true,
	//}
	conn, err := rabbitmq.Dial(cfg.RabbitMQ.Producer.DocumentNameSender.URL)
	if err != nil {
		logger.Fatal("connection dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("channel: %w", err)
	}
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := newRabbitMQProducer(cfg, ch); err != nil {
		logger.Fatal("rabbitmq producer sender: %w", zap.Error(err))
	}

	// Storage
	logger.Info("starting storage...")

	storage := psqlrep.New(pg)

	// Services
	logger.Info("starting services...")

	//documentImpl := documentsrvc.NewKafka(
	//	documentSenderProducer,
	//	storage,
	//	logger,
	//)
	documentImpl := documentsrvc.NewRabbiqMQ(
		documentsrvc.RabbitMQConfig{
			Exchange: cfg.RabbitMQ.Producer.DocumentNameSender.Exchange,
			Key:      cfg.RabbitMQ.Producer.DocumentNameSender.QueueName,
		},
		ch,
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

	fileImpl := filesrvc.New(
		storage,
		logger,
	)

	rabbitmqHandler := rabbitmqconsumer.New(
		ch,
		rabbitmqconsumer.Config{
			Queue: cfg.RabbitMQ.Consumer.DocumentSaver.QueueName,
		},
		logger,
		requestImpl,
	)
	if err := rabbitmqHandler.DocumentSaver(); err != nil {
		logger.Fatal("rabbitmq consumer sender: %w", zap.Error(err))
	}
	// HTTP server
	logger.Info("starting HTTP server...")

	app := registerHTPP(cfg, logger,
		documentImpl,
		requestImpl,
		s3Impl,
		fileImpl,
	)

	// Kafka consumers
	logger.Info("starting kafka consumer...")

	// if err = startKafkaConsumers(cfg.Kafka.KafkaConsumer, logger); err != nil {
	// 	logger.Fatal("kafka consumer: %w", zap.Error(err))
	// }

	if err = app.Listen(cfg.App.Address); err != nil {
		logger.Fatal("Ошибка запуска сервера: %v", zap.Error(err))
	}
}
