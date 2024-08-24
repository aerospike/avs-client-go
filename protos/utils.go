package protos

import (
	"fmt"
	"reflect"
)

func ConvertToKey(namespace string, set *string, key any) (*Key, error) {
	k := &Key{
		Namespace: namespace,
		Set:       set,
	}

	switch v := key.(type) {
	case string:
		k.Value = &Key_StringValue{StringValue: v}
	case []byte:
		k.Value = &Key_BytesValue{BytesValue: v}
	case int32:
		k.Value = &Key_IntValue{IntValue: v}
	case int64:
		k.Value = &Key_LongValue{LongValue: v}
	case int:
		k.Value = &Key_LongValue{LongValue: int64(v)}
	default:
		return nil, fmt.Errorf("unsupported key type: %T", key)
	}

	return k, nil
}

func ConvertFromKey(protoKey *Key) (namespace string, set *string, key any, err error) {
	namespace = protoKey.Namespace
	set = protoKey.Set

	switch v := protoKey.GetValue().(type) {
	case *Key_StringValue:
		key = v.StringValue
	case *Key_BytesValue:
		key = v.BytesValue
	case *Key_IntValue:
		key = v.IntValue
	case *Key_LongValue:
		key = v.LongValue
	default:
		err = fmt.Errorf("unsupported key value type: %T", v)
	}

	return
}

func ConvertToValue(value any) (*Value, error) {
	var convertedValue *Value

	switch v := value.(type) {
	case string:
		convertedValue = &Value{Value: &Value_StringValue{StringValue: v}}
	case []byte:
		convertedValue = &Value{Value: &Value_BytesValue{BytesValue: v}}
	case int32:
		convertedValue = &Value{Value: &Value_IntValue{IntValue: v}}
	case int64:
		convertedValue = &Value{Value: &Value_LongValue{LongValue: v}}
	case int:
		convertedValue = &Value{Value: &Value_LongValue{LongValue: int64(v)}}
	case float32:
		convertedValue = &Value{Value: &Value_FloatValue{FloatValue: v}}
	case float64:
		convertedValue = &Value{Value: &Value_DoubleValue{DoubleValue: v}}
	case *Map:
		convertedValue = &Value{Value: &Value_MapValue{MapValue: v}}
	case *List:
		convertedValue = &Value{Value: &Value_ListValue{ListValue: v}}
	case *Vector:
		convertedValue = &Value{Value: &Value_VectorValue{VectorValue: v}}
	case bool:
		convertedValue = &Value{Value: &Value_BooleanValue{BooleanValue: v}}
	case []float32:
		return ConvertToValue(CreateFloat32Vector(v))
	case []bool:
		return ConvertToValue(CreateBoolVector(v))
	default:
		val := reflect.ValueOf(value)
		switch val.Kind() { //nolint:exhaustive // default handles it
		case reflect.Map:
			m, err := ConvertToMapValue(val.Interface())
			if err != nil {
				return nil, err
			}

			return ConvertToValue(m)
		case reflect.Slice:
			l, err := ConvertToList(val.Interface())
			if err != nil {
				return nil, err
			}

			return ConvertToValue(l)
		default:
			return nil, fmt.Errorf("unsupported value type: %T", value)
		}
	}

	return convertedValue, nil
}

func ConvertToMapKey(key any) (*MapKey, error) {
	var keyValue *MapKey

	switch v := key.(type) {
	case string:
		keyValue = &MapKey{Value: &MapKey_StringValue{StringValue: v}}
	case []byte:
		keyValue = &MapKey{Value: &MapKey_BytesValue{BytesValue: v}}
	case int32:
		keyValue = &MapKey{Value: &MapKey_IntValue{IntValue: v}}
	case int64:
		keyValue = &MapKey{Value: &MapKey_LongValue{LongValue: v}}
	case int:
		keyValue = &MapKey{Value: &MapKey_LongValue{LongValue: int64(v)}}
	default:
		return nil, fmt.Errorf("unsupported key type: %T", key)
	}

	return keyValue, nil
}

func ConvertFromMapKey(key *MapKey) (any, error) {
	switch v := key.Value.(type) {
	case *MapKey_StringValue:
		return v.StringValue, nil
	case *MapKey_BytesValue:
		return v.BytesValue, nil
	case *MapKey_IntValue:
		return v.IntValue, nil
	case *MapKey_LongValue:
		return v.LongValue, nil
	default:
		return nil, fmt.Errorf("unsupported map key value type: %T", key.Value)
	}
}

func ConvertToMapValue(value any) (*Map, error) {
	val := reflect.ValueOf(value)
	entries := make([]*MapEntry, 0, val.Len())

	for _, e := range val.MapKeys() {
		key, err := ConvertToMapKey(e.Interface())
		if err != nil {
			return nil, fmt.Errorf("unsupported map key: %w", err)
		}

		val, err := ConvertToValue(val.MapIndex(e).Interface())
		if err != nil {
			return nil, fmt.Errorf("unsupported map value: %w", err)
		}

		entries = append(entries, &MapEntry{
			Key:   key,
			Value: val,
		})
	}

	return &Map{
		Entries: entries,
	}, nil
}

func ConvertFromMapValue(value *Map) (map[any]any, error) {
	entries := make(map[any]any, len(value.Entries))

	for _, entry := range value.Entries {
		k, err := ConvertFromMapKey(entry.GetKey())
		if err != nil {
			return nil, fmt.Errorf("unsupported map key: %w", err)
		}

		v, err := ConvertFromValue(entry.GetValue())
		if err != nil {
			return nil, fmt.Errorf("unsupported map value: %w", err)
		}

		entries[k] = v
	}

	return entries, nil
}

func ConvertToList(value any) (*List, error) {
	refValue := reflect.ValueOf(value)

	if refValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("unsupported list value: %T", value)
	}

	entries := make([]*Value, 0, refValue.Len())

	for i := 0; i < refValue.Len(); i++ {
		v, err := ConvertToValue(refValue.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("unsupported list value: %w", err)
		}

		entries = append(entries, v)
	}

	return &List{
		Entries: entries,
	}, nil
}

func ConvertFromListValue(value *List) ([]any, error) {
	entries := make([]any, len(value.Entries))

	for idx, entry := range value.Entries {
		val, err := ConvertFromValue(entry)
		if err != nil {
			return nil, fmt.Errorf("unsupported list value: %w", err)
		}

		entries[idx] = val
	}

	return entries, nil
}

func ConvertFromValue(value *Value) (any, error) {
	switch v := value.Value.(type) {
	case *Value_StringValue:
		return v.StringValue, nil
	case *Value_BytesValue:
		return v.BytesValue, nil
	case *Value_IntValue:
		return v.IntValue, nil
	case *Value_LongValue:
		return v.LongValue, nil
	case *Value_FloatValue:
		return v.FloatValue, nil
	case *Value_DoubleValue:
		return v.DoubleValue, nil
	case *Value_MapValue:
		return ConvertFromMapValue(v.MapValue)
	case *Value_ListValue:
		return ConvertFromListValue(v.ListValue)
	case *Value_VectorValue:
		return ConvertFromVector(v.VectorValue)
	case *Value_BooleanValue:
		return v.BooleanValue, nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T", value.Value)
	}
}

func ConvertToFields(recordData map[string]any) ([]*Field, error) {
	fields := make([]*Field, 0, len(recordData))

	for k, v := range recordData {
		convertedValue, err := ConvertToValue(v)
		if err != nil {
			return nil, fmt.Errorf("error converting field value for key '%s': %w", k, err)
		}

		fields = append(fields, &Field{
			Name:  k,
			Value: convertedValue,
		})
	}

	return fields, nil
}

func ConvertFromFields(fields []*Field) (map[string]any, error) {
	recordData := make(map[string]any, len(fields))

	for _, field := range fields {
		convertedValue, err := ConvertFromValue(field.GetValue())
		if err != nil {
			return nil, fmt.Errorf("error converting field value for name '%s': %w", field.GetName(), err)
		}

		recordData[field.GetName()] = convertedValue
	}

	return recordData, nil
}

func CreateFloat32Vector(vector []float32) *Vector {
	return &Vector{
		Data: &Vector_FloatData{
			FloatData: &FloatData{Value: vector},
		},
	}
}

func CreateBoolVector(vector []bool) *Vector {
	return &Vector{
		Data: &Vector_BoolData{
			BoolData: &BoolData{Value: vector},
		},
	}
}

func ConvertFromVector(vector *Vector) (any, error) {
	switch v := vector.Data.(type) {
	case *Vector_FloatData:
		return v.FloatData.Value, nil
	case *Vector_BoolData:
		return v.BoolData.Value, nil
	default:
		return nil, fmt.Errorf("unsupported vector data type: %T", vector.Data)
	}
}
