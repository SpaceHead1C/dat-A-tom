package rmq

import (
	"context"
	"datatom/internal/domain"
	pkgrmq "datatom/pkg/message_broker/rmq"
	"encoding/json"
	rmq "github.com/wagslane/go-rabbitmq"
)

func (b *Broker) SendRecord(ctx context.Context, req domain.SendRecordRequest) error {
	msg, err := json.Marshal(recordToSchema(req.Record))
	if err != nil {
		return err
	}
	p := pkgrmq.NewPublishing(
		b.publisher,
		req.RoutingKeys,
		msg,
		rmq.WithPublishOptionsPersistentDelivery,
		rmq.WithPublishOptionsExchange(req.Exchange),
		rmq.WithPublishOptionsContentType("application/json"),
		rmq.WithPublishOptionsType(domain.DeliveryTypeRecord),
		rmq.WithPublishOptionsAppID(req.TomID.String()),
	)
	return p.Publish(ctx)
}
