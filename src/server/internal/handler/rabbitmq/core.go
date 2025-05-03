package rabbitmq

import (
	requestsrvc "clusterlizer/internal/service/request"

	rabbitmq "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

type Handler struct {
	ch          *rabbitmq.Channel
	log         *zap.SugaredLogger
	cfg         Config
	requestSrvc requestsrvc.Service
}

type Config struct {
	Queue string
}

func New(
	ch *rabbitmq.Channel,
	cfg Config,
	log *zap.SugaredLogger,
	requestSrvc requestsrvc.Service,
) Handler {
	return Handler{
		ch:          ch,
		cfg:         cfg,
		log:         log,
		requestSrvc: requestSrvc,
	}
}
