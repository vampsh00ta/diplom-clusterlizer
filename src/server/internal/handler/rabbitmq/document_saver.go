package rabbitmq

import (
	"clusterlizer/internal/entity"
	requestsrvc "clusterlizer/internal/service/request"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
)

//type documentSaverParams struct {
//	ID     string         `json:"id"`
//	Groups []entity.Group `json:"groups"`
//}

func (h Handler) DocumentSaver() error {
	ctx := context.Background()
	errorChan := make(chan error, 1)
	go func() {
		for err := range errorChan {
			h.log.Error("error occurred", zap.Error(err))
		}
	}()
	_ = ctx
	q, err := h.ch.QueueDeclare(
		h.cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare q: %w", err)
	}
	msgs, err := h.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare consumer: %w", err)
	}
	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.MessageId)

			input, err := h.decodeDocumentSaverMsg(msg.Body)
			if err != nil {
				errorChan <- err
				continue
			}
			err = h.requestSrvc.SaveResult(ctx, requestsrvc.SaveResultParams{
				ID:    entity.RequestID(input.ID),
				Graph: input.Graph,
			})
			if err != nil {
				errorChan <- err
				continue
			}
		}
	}()

	return nil
}
func (h Handler) decodeDocumentSaverMsg(b []byte) (entity.DocumentSaverReq, error) {
	var input entity.DocumentSaverReq
	if err := json.Unmarshal(b, &input); err != nil {
		return entity.DocumentSaverReq{}, err
	}

	if input.ID == "" {
		return entity.DocumentSaverReq{}, fmt.Errorf("nil id")
	}
	return input, nil

}
