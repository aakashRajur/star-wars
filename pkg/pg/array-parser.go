package pg

import (
	"github.com/jackc/pgx/pgtype"
)

func GetInterfaceFromValue(val pgtype.Value) interface{} {
	return val.Get()
}

func SafeParseValue(value interface{}) interface{} {
	parsed := make([]interface{}, 0)

	switch value.(type) {
	case *pgtype.BoolArray:
		{
			ptr := value.(*pgtype.BoolArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.Int4Array:
		{
			ptr := value.(*pgtype.Int4Array)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.Int8Array:
		{
			ptr := value.(*pgtype.Int8Array)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.Float4Array:
		{
			ptr := value.(*pgtype.Float4Array)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.Float8Array:
		{
			ptr := value.(*pgtype.Float8Array)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.DateArray:
		{
			ptr := value.(*pgtype.DateArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.TimestamptzArray:
		{
			ptr := value.(*pgtype.TimestamptzArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.TimestampArray:
		{
			ptr := value.(*pgtype.TimestampArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.EnumArray:
		{
			ptr := value.(*pgtype.EnumArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.TextArray:
		{
			ptr := value.(*pgtype.TextArray)
			for _, each := range ptr.Elements {
				parsed = append(parsed, GetInterfaceFromValue(&each))
			}
			return parsed
		}
	case *pgtype.UntypedTextArray:
		{
			ptr := value.(*pgtype.UntypedTextArray)
			parsed = append(parsed, ptr.Elements)
		}
	}
	return value
}
