package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ValueRepository interface {
	SetValue(context.Context, SetValueRequest) (*Value, error)
}

type Value struct {
	RecordID   uuid.UUID
	PropertyID uuid.UUID
	Type       Type
	RefTypeID  uuid.UUID
	Value      any
	Sum        string
	ChangeAt   time.Time
}

type SetValueRequest struct {
	RecordID   uuid.UUID
	PropertyID uuid.UUID
	Type       Type
	RefTypeID  uuid.UUID
	Value      any
}

type ValueJSONSchema struct {
	V any `json:"v"`
}

func ValueAsJSON(v any, t Type) ([]byte, error) {
	switch t {
	case TypeText:
		switch v.(type) {
		case string:
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeNumber:
		switch v.(type) {
		case int, float64:
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeBool:
		switch v.(type) {
		case bool:
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeDate:
		switch v.(type) {
		case time.Time:
		case string:
			if _, err := time.Parse(time.RFC3339, v.(string)); err != nil {
				return nil, fmt.Errorf("%w date error: %s", ErrParseError, err)
			}
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeUUID, TypeReference:
		switch v.(type) {
		case uuid.UUID:
		case string:
			if _, err := uuid.Parse(v.(string)); err != nil {
				return nil, fmt.Errorf("%w UUID error: %s", ErrParseError, err)
			}
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	default:
		return nil, fmt.Errorf("%w %s", ErrUnexpectedTypePG, t.String())
	}
	return json.Marshal(ValueJSONSchema{v})
}

func ValidatedValue(v any, t Type) (any, error) {
	var out any
	switch t {
	case TypeText:
		switch v.(type) {
		case string:
			out = v
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeNumber:
		switch v.(type) {
		case int, float64:
			out = v
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeBool:
		switch v.(type) {
		case bool:
			out = v
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeDate:
		switch v.(type) {
		case time.Time:
			out = v
		case string:
			x, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return nil, fmt.Errorf("%w date error: %s", ErrParseError, err)
			}
			out = x
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	case TypeUUID, TypeReference:
		switch v.(type) {
		case uuid.UUID:
			out = v
		case string:
			x, err := uuid.Parse(v.(string))
			if err != nil {
				return nil, fmt.Errorf("%w UUID error: %s", ErrParseError, err)
			}
			out = x
		default:
			return nil, fmt.Errorf("%w %T for %s", ErrUnexpectedTypePG, v, t.String())
		}
	default:
		return nil, fmt.Errorf("%w %s", ErrUnexpectedTypePG, t.String())
	}
	return out, nil
}