package domain

type Type uint

const (
	UndefinedType = iota
	TypeNumber
	TypeText
	TypeBool
	TypeDate
	TypeUUID
	TypeReference
)

func (tp Type) String() string {
	switch tp {
	case TypeNumber:
		return "number"
	case TypeText:
		return "text"
	case TypeBool:
		return "boolean"
	case TypeDate:
		return "date"
	case TypeUUID:
		return "UUID"
	case TypeReference:
		return "reference"
	default:
		return "undefined"
	}
}

func (tp Type) Code() string {
	switch tp {
	case TypeNumber:
		return "number"
	case TypeText:
		return "text"
	case TypeBool:
		return "bool"
	case TypeDate:
		return "date"
	case TypeUUID:
		return "uuid"
	case TypeReference:
		return "ref"
	default:
		return "undefined"
	}
}

func TypesToCodes(ts []Type) []string {
	if ts == nil {
		return nil
	}
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		out = append(out, t.Code())
	}
	return out
}
