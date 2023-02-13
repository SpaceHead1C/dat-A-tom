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
	for _, tp := range ts {
		out = append(out, tp.Code())
	}
	return out
}

func TypeFromCode(code string) Type {
	switch code {
	case "number":
		return TypeNumber
	case "text":
		return TypeText
	case "bool":
		return TypeBool
	case "date":
		return TypeDate
	case "uuid":
		return TypeUUID
	case "ref":
		return TypeReference
	default:
		return UndefinedType
	}
}

func TypesFromCodes(codes []string) []Type {
	if codes == nil {
		return nil
	}
	out := make([]Type, 0, len(codes))
	for _, code := range codes {
		out = append(out, TypeFromCode(code))
	}
	return out
}
