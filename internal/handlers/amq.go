package handlers

import (
	"context"
	"datatom/internal/api"
	"datatom/pkg/log"
	"encoding/json"
	"fmt"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"time"
)

const deliveryTypeValue = "value"

type ConsumeHandlerConfig struct {
	Logger       *zap.SugaredLogger
	Timeout      time.Duration
	ValueManager *api.ValueManager
}

func NewConsumeHandler(c ConsumeHandlerConfig) rmq.Handler {
	return func(d rmq.Delivery) (action rmq.Action) {
		if c.Timeout == 0 {
			c.Timeout = time.Second * 2
		}
		if c.Logger == nil {
			c.Logger = log.GlobalLogger()
		}
		ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
		defer cancel()
		var err error
		switch d.Type {
		case deliveryTypeValue:
			err = processMessageWithValue(ctx, c.ValueManager, d.Body)
		default:
			err = fmt.Errorf("unexpected delivery type %s", d.Type)
		}
		if err != nil {
			action = rmq.NackDiscard
			c.Logger.Errorln(err)
		}
		return action
	}
}

func processMessageWithValue(ctx context.Context, man *api.ValueManager, message []byte) error {
	var schema SetValueRequestSchema
	if err := json.Unmarshal(message, &schema); err != nil {
		return err
	}
	req, err := schema.SetValueRequest()
	if err != nil {
		return err
	}
	if _, err := man.Set(ctx, req); err != nil {
		return err
	}
	return nil
}
