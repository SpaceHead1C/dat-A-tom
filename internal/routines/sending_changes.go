package routines

import (
	"context"
	"datatom/internal/api"
	"datatom/internal/domain"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type SendChangedDataConfig struct {
	Logger               *zap.SugaredLogger
	ReferenceTypeManager *api.RefTypeManager
	RecordManager        *api.RecordManager
	PropertyManager      *api.PropertyManager
	ValueManager         *api.ValueManager
	ChangedDataManager   *api.ChangedDataManager
	StoredConfigsManager *api.StoredConfigsManager
	DBManager            *api.DBManager
	Exchange             string
	RoutingKeys          []string
}

func sendChangedData(c SendChangedDataConfig) error {
	var tomID uuid.UUID
	storedValue, err := c.StoredConfigsManager.Get(context.Background(), domain.StoredConfigTomID)
	if err != nil {
		if errors.Is(err, domain.ErrStoredConfigTomIDNotSet) {
			return fmt.Errorf("changed data sending failed cause error: %w", err)
		}
		return fmt.Errorf("get tom ID for sending changed data error: %w", err)
	}
	if err := storedValue.ScanStoredConfigValue(&tomID); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	changes, err := c.ChangedDataManager.Get(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("get changed data error: %w", err)
	}
	last, errOut := processChanges(c, tomID, changes)
	if errOut != nil {
		last--
	}
	if last >= 0 {
		if err := c.ChangedDataManager.Purge(context.Background(), changes[last].ID, nil); err != nil {
			c.Logger.Errorf("clear registration of sent changes error: %s", err)
		}
	}
	return errOut
}

func processChanges(c SendChangedDataConfig, tomID uuid.UUID, changes []domain.ChangedData) (int, error) {
	for i, change := range changes {
		sender, err := newSender(c, tomID, change)
		if err != nil {
			return i, err
		}
		if err := executeSender(c.DBManager, sender); err != nil {
			return i, err
		}
	}
	return len(changes) - 1, nil
}

func newSender(c SendChangedDataConfig, tomID uuid.UUID, change domain.ChangedData) (api.Sender, error) {
	switch change.DataType {
	case domain.ChangedDataValue:
		value, err := c.ValueManager.GetByKey(context.Background(), change.Key)
		if err != nil {
			return nil, fmt.Errorf("get changed value error: %s", err)
		}
		return c.ValueManager.GetSender(domain.SendValueRequest{
			Value:       *value,
			TomID:       tomID,
			Exchange:    c.Exchange,
			RoutingKeys: c.RoutingKeys,
		}), nil
	case domain.ChangedDataRecord:
		record, err := c.RecordManager.GetByKey(context.Background(), change.Key)
		if err != nil {
			return nil, fmt.Errorf("get changed record error: %s", err)
		}
		return c.RecordManager.GetSender(domain.SendRecordRequest{
			Record:      *record,
			TomID:       tomID,
			Exchange:    c.Exchange,
			RoutingKeys: c.RoutingKeys,
		}), nil
	case domain.ChangedDataProperty:
		property, err := c.PropertyManager.GetByKey(context.Background(), change.Key)
		if err != nil {
			return nil, fmt.Errorf("get changed property error: %s", err)
		}
		return c.PropertyManager.GetSender(domain.SendPropertyRequest{
			Property:    *property,
			TomID:       tomID,
			Exchange:    c.Exchange,
			RoutingKeys: c.RoutingKeys,
		}), nil
	case domain.ChangedDataRefType:
		refType, err := c.ReferenceTypeManager.GetByKey(context.Background(), change.Key)
		if err != nil {
			return nil, fmt.Errorf("get changed reference type error: %s", err)
		}
		return c.ReferenceTypeManager.GetSender(domain.SendRefTypeRequest{
			RefType:     *refType,
			TomID:       tomID,
			Exchange:    c.Exchange,
			RoutingKeys: c.RoutingKeys,
		}), nil
	default:
		return nil, fmt.Errorf("%w \"%s\" of changed data", domain.ErrUnknownType, change.DataType.String())
	}
}

func executeSender(man *api.DBManager, sender api.Sender) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tx, err := man.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction error: %s", err)
	}
	equals, err := sender.SumEqualsSent(ctx, tx)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("get sent data error: %s", err)
	}
	if equals {
		tx.Rollback(context.Background())
		return nil
	}
	if err := sender.SetSentState(ctx, tx); err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("set sent state error: %s", err)
	}
	if err := sender.Send(ctx); err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("publish error: %s", err)
	}
	tx.Commit(context.Background())
	return nil
}
