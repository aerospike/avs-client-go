package protos

func ConvertToKey(key any) (keyValue isKey_Value) {
	switch v := key.(type) {
	case string:
		keyValue = &Key_StringValue{StringValue: v}
	case []byte:
		keyValue = &Key_BytesValue{BytesValue: v}
	case int32:
		keyValue = &Key_IntValue{IntValue: v}
	case int64:
		keyValue = &Key_LongValue{LongValue: v}
	default:
		// handle unsupported types or error
	}
	return keyValue
}

func ConvertToValue(value any) *Value {
	switch v := value.(type) {
	case string:
		return &Value{Value: &Value_StringValue{StringValue: v}}
	case []byte:
		return &Value{Value: &Value_BytesValue{BytesValue: v}}
	case int32:
		return &Value{Value: &Value_IntValue{IntValue: v}}
	case int64:
		return &Value{Value: &Value_LongValue{LongValue: v}}
	case float32:
		return &Value{Value: &Value_FloatValue{FloatValue: v}}
	case float64:
		return &Value{Value: &Value_DoubleValue{DoubleValue: v}}
	case *Map:
		return &Value{Value: &Value_MapValue{MapValue: v}}
	case *List:
		return &Value{Value: &Value_ListValue{ListValue: v}}
	case *Vector:
		return &Value{Value: &Value_VectorValue{VectorValue: v}}
	case bool:
		return &Value{Value: &Value_BooleanValue{BooleanValue: v}}
	default:
		return nil
	}
}

func ConvertFromValue(value *Value) any {
	switch v := value.Value.(type) {
	case *Value_StringValue:
		return v.StringValue
	case *Value_BytesValue:
		return v.BytesValue
	case *Value_IntValue:
		return v.IntValue
	case *Value_LongValue:
		return v.LongValue
	case *Value_FloatValue:
		return v.FloatValue
	case *Value_DoubleValue:
		return v.DoubleValue
	case *Value_MapValue:
		return v.MapValue
	case *Value_ListValue:
		return v.ListValue
	case *Value_VectorValue:
		return v.VectorValue
	case *Value_BooleanValue:
		return v.BooleanValue
	default:
		return nil
	}
}

func ConvertToFields(recordData map[string]any) []*Field {
	fields := make([]*Field, 0, len(recordData))
	for k, v := range recordData {
		fields = append(fields, &Field{
			Name:  k,
			Value: ConvertToValue(v),
		})
	}
	return fields
}

func ConvertFromFields(fields []*Field) map[string]any {
	recordData := make(map[string]any, len(fields))
	for _, field := range fields {
		recordData[field.GetName()] = ConvertFromValue(field.GetValue())
	}
	return recordData
}

func CreateFloat32Vector(vector []float32) *Vector {
	return &Vector{
		Data: &Vector_FloatData{FloatData: &FloatData{Value: vector}},
	}
}

func CreateBoolVector(vector []bool) *Vector {
	return &Vector{
		Data: &Vector_BoolData{BoolData: &BoolData{Value: vector}},
	}
}
