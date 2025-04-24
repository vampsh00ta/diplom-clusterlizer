package rabbitmq

import (
	"clusterlizer/internal/entity"
	requestsrvc "clusterlizer/internal/service/request"
	"clusterlizer/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
)

type documentSaverParams struct {
	ID     string         `json:"id"`
	Groups []entity.Group `json:"groups"`
}

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

			bytes, ID, err := h.decodeDocumentSaverMsg(msg.Body)
			if err != nil {
				errorChan <- err
				continue
			}
			_, err = h.requestSrvc.UpdateRequest(ctx, requestsrvc.UpdateRequestParams{
				ID:     entity.RequestID(ID),
				Result: utils.NewOptional(&bytes),
				Status: utils.NewOptional(entity.StatusDone),
			})
			if err != nil {
				errorChan <- err
				continue
			}
		}
	}()

	return nil
}
func (h Handler) decodeDocumentSaverMsg(b []byte) ([]byte, string, error) {
	var input documentSaverParams
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, "", err
	}
	if input.ID == "" {
		return nil, "", fmt.Errorf("nil id")
	}
	var res entity.Groups
	res = input.Groups

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, "", err
	}

	return bytes, input.ID, nil

}
