package handlers

import (
	"context"
	"datatom/internal/api"
	"datatom/internal/domain"
	"datatom/pkg/log"
	"encoding/json"
	"errors"
	"fmt"
	rmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"time"
)

const (
	deliveryTypeValue    = "value"
	deliveryTypeProperty = "property"
	deliveryTypeRecord   = "record"
)

type ConsumeHandlerConfig struct {
	Logger          *zap.SugaredLogger
	Timeout         time.Duration
	ValueManager    *api.ValueManager
	PropertyManager *api.PropertyManager
	RecordManager   *api.RecordManager
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
		var isInnerError bool
		switch d.Type {
		case deliveryTypeValue:
			isInnerError, err = processMessageWithValue(ctx, c.ValueManager, d.Body)
		case deliveryTypeProperty:
			isInnerError, err = processMessageWithProperty(ctx, c.PropertyManager, d.Body)
		case deliveryTypeRecord:
			isInnerError, err = processMessageWithRecord(ctx, c.RecordManager, d.Body)
		default:
			err = fmt.Errorf("unexpected delivery type %s", d.Type)
		}
		if err != nil {
			action = rmq.NackDiscard
			template := "process message %s error: %s"
			if isInnerError {
				c.Logger.Errorf(template, d.MessageId, err)
			} else {
				c.Logger.Infof(template, d.MessageId, err)
			}
		}
		return action
	}
}

func processMessageWithValue(ctx context.Context, man *api.ValueManager, message []byte) (bool, error) {
	var schema SetValueRequestSchema
	if err := json.Unmarshal(message, &schema); err != nil {
		return false, err
	}
	req, err := schema.SetValueRequest()
	if err != nil {
		return false, err
	}
	if _, err := man.Set(ctx, req); err != nil {
		return !errors.Is(err, domain.ErrNotFound), err
	}
	return false, nil
}

func processMessageWithProperty(ctx context.Context, man *api.PropertyManager, message []byte) (bool, error) {
	var schema UpdPropertyRequestSchema
	if err := json.Unmarshal(message, &schema); err != nil {
		return false, err
	}
	req, err := schema.UpdPropertyRequest()
	if err != nil {
		return false, err
	}
	if _, err := man.Update(ctx, req); err != nil {
		return !errors.Is(err, domain.ErrNotFound), err
	}
	return false, nil
}

func processMessageWithRecord(ctx context.Context, man *api.RecordManager, message []byte) (bool, error) {
	var schema UpdRecordRequestSchema
	if err := json.Unmarshal(message, &schema); err != nil {
		return false, err
	}
	req, err := schema.UpdRecordRequest()
	if err != nil {
		return false, err
	}
	if _, err := man.Update(ctx, req); err != nil {
		return !errors.Is(err, domain.ErrNotFound), err
	}
	return false, nil
}
