package protos

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Ptr(str string) *string {
	ptr := str
	return &ptr
}
func TestConvertToKey(t *testing.T) {
	testCases := []struct {
		input       any
		expected    *Key
		expectedErr error
	}{
		{
			input: "testString",
			expected: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_StringValue{
					StringValue: "testString",
				},
			},
		},
		{
			input: []byte{0x01, 0x02, 0x03},
			expected: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
		},
		{
			input: int32(123),
			expected: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_IntValue{
					IntValue: 123,
				},
			},
		},
		{
			input: int64(123456789),
			expected: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_LongValue{
					LongValue: 123456789,
				},
			},
		},
		{
			input: int(123456789),
			expected: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_LongValue{
					LongValue: 123456789,
				},
			},
		},
		{
			input:       &keyUnknown{},
			expected:    nil, // Unsupported type
			expectedErr: fmt.Errorf("unsupported key type: *protos.keyUnknown"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertToKey("testNamespace", Ptr("testSet"), tc.input)

		assert.Equal(t, tc.expected, result)
		assert.Equal(t, tc.expectedErr, err)
	}
}

type keyUnknown struct{}

func (*keyUnknown) isKey_Value() {} //nolint:revive,stylecheck // Grpc generated

func TestConvertFromKey(t *testing.T) {
	testCases := []struct {
		input             *Key
		expectedNamespace string
		expectedSet       *string
		expectedKey       any
		expectedErr       error
	}{
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_StringValue{
					StringValue: "testString",
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       Ptr("testSet"),
			expectedKey:       "testString",
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       Ptr("testSet"),
			expectedKey:       []byte{0x01, 0x02, 0x03},
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_IntValue{
					IntValue: 123,
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       Ptr("testSet"),
			expectedKey:       int32(123),
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value: &Key_LongValue{
					LongValue: 123456789,
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       Ptr("testSet"),
			expectedKey:       int64(123456789),
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       Ptr("testSet"),
				Value:     &keyUnknown{},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       Ptr("testSet"),
			expectedKey:       nil, // Unsupported or nil input
			expectedErr:       fmt.Errorf("unsupported key value type: *protos.keyUnknown"),
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resultNamespace, resultSet, resultKey, err := ConvertFromKey(tc.input)

			assert.Equal(t, tc.expectedNamespace, resultNamespace)
			assert.Equal(t, tc.expectedSet, resultSet)
			assert.Equal(t, tc.expectedKey, resultKey)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestConvertToValue(t *testing.T) {
	testCases := []struct {
		input       any
		expected    *Value
		expectedErr bool
	}{
		{
			input: "testString",
			expected: &Value{
				Value: &Value_StringValue{
					StringValue: "testString",
				},
			},
		},
		{
			input: []byte{0x01, 0x02, 0x03},
			expected: &Value{
				Value: &Value_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
		},
		{
			input: int32(123),
			expected: &Value{
				Value: &Value_IntValue{
					IntValue: 123,
				},
			},
		},
		{
			input: int64(123456789),
			expected: &Value{
				Value: &Value_LongValue{
					LongValue: 123456789,
				},
			},
		},
		{
			input: int(456),
			expected: &Value{
				Value: &Value_LongValue{
					LongValue: int64(456),
				},
			},
		},
		{
			input: float32(123.45),
			expected: &Value{
				Value: &Value_FloatValue{
					FloatValue: 123.45,
				},
			},
		},
		{
			input: float64(123.456789),
			expected: &Value{
				Value: &Value_DoubleValue{
					DoubleValue: 123.456789,
				},
			},
		},
		{
			input: true,
			expected: &Value{
				Value: &Value_BooleanValue{
					BooleanValue: true,
				},
			},
		},
		{
			input: map[any]any{"key": "value"},
			expected: &Value{
				Value: &Value_MapValue{
					MapValue: &Map{
						Entries: []*MapEntry{
							{
								Key: &MapKey{
									Value: &MapKey_StringValue{
										StringValue: "key",
									},
								},
								Value: &Value{
									Value: &Value_StringValue{
										StringValue: "value",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: map[string]string{"key": "value"},
			expected: &Value{
				Value: &Value_MapValue{
					MapValue: &Map{
						Entries: []*MapEntry{
							{
								Key: &MapKey{
									Value: &MapKey_StringValue{
										StringValue: "key",
									},
								},
								Value: &Value{
									Value: &Value_StringValue{
										StringValue: "value",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: map[int]float64{10: 3.124},
			expected: &Value{
				Value: &Value_MapValue{
					MapValue: &Map{
						Entries: []*MapEntry{
							{
								Key: &MapKey{
									Value: &MapKey_LongValue{
										LongValue: int64(10),
									},
								},
								Value: &Value{
									Value: &Value_DoubleValue{
										DoubleValue: 3.124,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: []any{"item1", "item2"},
			expected: &Value{
				Value: &Value_ListValue{
					ListValue: &List{
						Entries: []*Value{
							{
								Value: &Value_StringValue{
									StringValue: "item1",
								},
							},
							{
								Value: &Value_StringValue{
									StringValue: "item2",
								},
							},
						},
					},
				},
			},
		},
		{
			input: []string{"item1", "item2"},
			expected: &Value{
				Value: &Value_ListValue{
					ListValue: &List{
						Entries: []*Value{
							{
								Value: &Value_StringValue{
									StringValue: "item1",
								},
							},
							{
								Value: &Value_StringValue{
									StringValue: "item2",
								},
							},
						},
					},
				},
			},
		},
		{
			input: []int{1, 2},
			expected: &Value{
				Value: &Value_ListValue{
					ListValue: &List{
						Entries: []*Value{
							{
								Value: &Value_LongValue{
									LongValue: int64(1),
								},
							},
							{
								Value: &Value_LongValue{
									LongValue: int64(2),
								},
							},
						},
					},
				},
			},
		},
		{
			input: []float32{1, 2},
			expected: &Value{
				Value: &Value_VectorValue{
					VectorValue: &Vector{
						Data: &Vector_FloatData{
							FloatData: &FloatData{
								Value: []float32{1, 2},
							},
						},
					},
				},
			},
		},
		{
			input: []bool{true, false},
			expected: &Value{
				Value: &Value_VectorValue{
					VectorValue: &Vector{
						Data: &Vector_BoolData{
							BoolData: &BoolData{
								Value: []bool{true, false},
							},
						},
					},
				},
			},
		},
		{
			input:       map[int]any{10: struct{}{}},
			expected:    nil,
			expectedErr: true,
		},
		{
			input:       []any{struct{}{}},
			expected:    nil,
			expectedErr: true,
		},
		{
			input:       struct{}{}, // Unsupported type
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := ConvertToValue(tc.input)

			assert.Equal(t, tc.expected, result)

			if tc.expectedErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestConvertToMapKey(t *testing.T) {
	testCases := []struct {
		input       any
		expected    *MapKey
		expectedErr error
	}{
		{
			input: "testString",
			expected: &MapKey{
				Value: &MapKey_StringValue{
					StringValue: "testString",
				},
			},
		},
		{
			input: int32(123),
			expected: &MapKey{
				Value: &MapKey_IntValue{
					IntValue: 123,
				},
			},
		},
		{
			input: int64(123456789),
			expected: &MapKey{
				Value: &MapKey_LongValue{
					LongValue: 123456789,
				},
			},
		},
		{
			input: int(123456789),
			expected: &MapKey{
				Value: &MapKey_LongValue{
					LongValue: 123456789,
				},
			},
		},
		{
			input: []byte{0x01, 0x02, 0x03},
			expected: &MapKey{
				Value: &MapKey_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
		},
		{
			input:       struct{}{},
			expected:    nil, // Unsupported type
			expectedErr: fmt.Errorf("unsupported key type: struct {}"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertToMapKey(tc.input)

		assert.Equal(t, tc.expected, result)
		assert.Equal(t, tc.expectedErr, err)
	}
}

type mapKeyValueUnknown struct{}

func (*mapKeyValueUnknown) isMapKey_Value() {} //nolint:revive,stylecheck // Grpc generated

func TestConvertFromMapKey(t *testing.T) {
	testCases := []struct {
		input       *MapKey
		expected    any
		expectedErr error
	}{
		{
			input: &MapKey{
				Value: &MapKey_StringValue{
					StringValue: "testString",
				},
			},
			expected: "testString",
		},
		{
			input: &MapKey{
				Value: &MapKey_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
			expected: []byte{0x01, 0x02, 0x03},
		},
		{
			input: &MapKey{
				Value: &MapKey_IntValue{
					IntValue: 123,
				},
			},
			expected: int32(123),
		},
		{
			input: &MapKey{
				Value: &MapKey_LongValue{
					LongValue: 123456789,
				},
			},
			expected: int64(123456789),
		},
		{
			input: &MapKey{
				Value: &mapKeyValueUnknown{},
			},
			expected:    nil,
			expectedErr: fmt.Errorf("unsupported map key value type: *protos.mapKeyValueUnknown"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertFromMapKey(tc.input)

		assert.Equal(t, tc.expected, result)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestConvertToMapValue(t *testing.T) {
	testCases := []struct {
		input       any
		expected    *Map
		expectedErr *string
	}{
		{
			input: map[any]any{"key": "value"},
			expected: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_StringValue{
								StringValue: "key",
							},
						},
						Value: &Value{
							Value: &Value_StringValue{
								StringValue: "value",
							},
						},
					},
				},
			},
		},
		{
			input: map[string]string{"key": "value"},
			expected: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_StringValue{
								StringValue: "key",
							},
						},
						Value: &Value{
							Value: &Value_StringValue{
								StringValue: "value",
							},
						},
					},
				},
			},
		},
		{
			input: map[int]float64{10: 3.124},
			expected: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_LongValue{
								LongValue: int64(10),
							},
						},
						Value: &Value{
							Value: &Value_DoubleValue{
								DoubleValue: 3.124,
							},
						},
					},
				},
			},
		},
		{
			input:       map[int]any{10: struct{}{}},
			expected:    nil,
			expectedErr: Ptr("unsupported map value: unsupported value type: struct {}"),
		},
		{
			input:       map[any]any{struct{}{}: 10},
			expected:    nil,
			expectedErr: Ptr("unsupported map key: unsupported key type: struct {}"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertToMapValue(tc.input)

		assert.Equal(t, tc.expected, result)

		if tc.expectedErr != nil {
			assert.ErrorContains(t, err, *tc.expectedErr)
		}
	}
}

func TestConvertToList(t *testing.T) {
	testCases := []struct {
		input       []any
		expected    *List
		expectedErr *string
	}{
		{
			input: []any{"item1", "item2"},
			expected: &List{
				Entries: []*Value{
					{
						Value: &Value_StringValue{
							StringValue: "item1",
						},
					},
					{
						Value: &Value_StringValue{
							StringValue: "item2",
						},
					},
				},
			},
		},
		{
			input: []any{1, 2},
			expected: &List{
				Entries: []*Value{
					{
						Value: &Value_LongValue{
							LongValue: int64(1),
						},
					},
					{
						Value: &Value_LongValue{
							LongValue: int64(2),
						},
					},
				},
			},
		},
		{
			input: []any{true, false},
			expected: &List{
				Entries: []*Value{
					{
						Value: &Value_BooleanValue{
							BooleanValue: true,
						},
					},
					{
						Value: &Value_BooleanValue{
							BooleanValue: false,
						},
					},
				},
			},
		},
		{
			input:       []any{struct{}{}},
			expected:    nil,
			expectedErr: Ptr("unsupported list value: unsupported value type: struct {}"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertToList(tc.input)

		assert.Equal(t, tc.expected, result)

		if tc.expectedErr != nil {
			assert.ErrorContains(t, err, *tc.expectedErr)
		}
	}
}

func TestConvertFromListValue(t *testing.T) {
	testCases := []struct {
		input       *List
		expected    []any
		expectedErr *string
	}{
		{
			input: &List{
				Entries: []*Value{
					{
						Value: &Value_StringValue{
							StringValue: "item1",
						},
					},
					{
						Value: &Value_StringValue{
							StringValue: "item2",
						},
					},
				},
			},
			expected: []any{"item1", "item2"},
		},
		{
			input: &List{
				Entries: []*Value{
					{
						Value: &Value_LongValue{
							LongValue: int64(1),
						},
					},
					{
						Value: &Value_LongValue{
							LongValue: int64(2),
						},
					},
				},
			},
			expected: []any{int64(1), int64(2)},
		},
		{
			input: &List{
				Entries: []*Value{
					{
						Value: &Value_BooleanValue{
							BooleanValue: true,
						},
					},
					{
						Value: &Value_BooleanValue{
							BooleanValue: false,
						},
					},
				},
			},
			expected: []any{true, false},
		},
		{
			input: &List{
				Entries: []*Value{
					{
						Value: &valueUnknown{},
					},
				},
			},
			expected:    nil,
			expectedErr: Ptr("unsupported list value: unsupported value type: *protos.valueUnknown"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertFromListValue(tc.input)

		assert.Equal(t, tc.expected, result)

		if tc.expectedErr != nil {
			assert.ErrorContains(t, err, *tc.expectedErr)
		}
	}
}

func TestConvertFromMapValue(t *testing.T) {
	var nilMap map[any]any
	testCases := []struct {
		input       *Map
		expected    any
		expectedErr *string
	}{
		{
			input: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_StringValue{
								StringValue: "key",
							},
						},
						Value: &Value{
							Value: &Value_StringValue{
								StringValue: "value",
							},
						},
					},
				},
			},
			expected: map[any]any{"key": "value"},
		},
		{
			input: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_StringValue{
								StringValue: "key",
							},
						},
						Value: &Value{
							Value: &Value_StringValue{
								StringValue: "value",
							},
						},
					},
				},
			},
			expected: map[any]any{"key": "value"},
		},
		{
			input: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_LongValue{
								LongValue: int64(10),
							},
						},
						Value: &Value{
							Value: &Value_DoubleValue{
								DoubleValue: 3.124,
							},
						},
					},
				},
			},
			expected: map[any]any{int64(10): 3.124},
		},
		{
			input: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &mapKeyValueUnknown{},
						},
						Value: &Value{
							Value: &Value_DoubleValue{
								DoubleValue: 3.124,
							},
						},
					},
				},
			},
			expected:    nilMap,
			expectedErr: Ptr("unsupported map key value type: *protos.mapKeyValueUnknown"),
		},
		{
			input: &Map{
				Entries: []*MapEntry{
					{
						Key: &MapKey{
							Value: &MapKey_LongValue{
								LongValue: int64(10),
							},
						},
						Value: &Value{
							Value: &valueUnknown{},
						},
					},
				},
			},
			expected:    nilMap,
			expectedErr: Ptr("unsupported map value: unsupported value type: *protos.valueUnknown"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertFromMapValue(tc.input)

		assert.Equal(t, tc.expected, result)

		if tc.expectedErr != nil {
			assert.ErrorContains(t, err, *tc.expectedErr)
		}
	}
}

type valueUnknown struct{}

func (*valueUnknown) isValue_Value() {} //nolint:revive,stylecheck // Grpc generated

func TestConvertFromValue(t *testing.T) {
	testCases := []struct {
		input       *Value
		expected    any
		expectedErr error
	}{
		{
			input: &Value{
				Value: &Value_StringValue{
					StringValue: "testString",
				},
			},
			expected: "testString",
		},
		{
			input: &Value{
				Value: &Value_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
			expected: []byte{0x01, 0x02, 0x03},
		},
		{
			input: &Value{
				Value: &Value_IntValue{
					IntValue: 123,
				},
			},
			expected: int32(123),
		},
		{
			input: &Value{
				Value: &Value_LongValue{
					LongValue: 123456789,
				},
			},
			expected: int64(123456789),
		},
		{
			input: &Value{
				Value: &Value_FloatValue{
					FloatValue: 123.45,
				},
			},
			expected: float32(123.45),
		},
		{
			input: &Value{
				Value: &Value_DoubleValue{
					DoubleValue: 123.456789,
				},
			},
			expected: float64(123.456789),
		},
		{
			input: &Value{
				Value: &Value_BooleanValue{
					BooleanValue: true,
				},
			},
			expected: true,
		},
		{
			input: &Value{
				Value: &Value_MapValue{
					MapValue: &Map{
						Entries: []*MapEntry{
							{
								Key: &MapKey{
									Value: &MapKey_StringValue{
										StringValue: "key",
									},
								},
								Value: &Value{
									Value: &Value_StringValue{
										StringValue: "value",
									},
								},
							},
						},
					},
				},
			},
			expected: map[any]any{"key": "value"},
		},
		{
			input: &Value{
				Value: &Value_ListValue{
					ListValue: &List{
						Entries: []*Value{
							{
								Value: &Value_StringValue{
									StringValue: "item1",
								},
							},
							{
								Value: &Value_StringValue{
									StringValue: "item2",
								},
							},
						},
					},
				},
			},
			expected: []any{"item1", "item2"},
		},
		{
			input: &Value{
				Value: &Value_VectorValue{
					VectorValue: &Vector{
						Data: &Vector_FloatData{
							FloatData: &FloatData{
								Value: []float32{1, 2},
							},
						},
					},
				},
			},
			expected: []float32{float32(1), float32(2)},
		},
		{
			input: &Value{
				Value: &valueUnknown{},
			},
			expected:    nil,
			expectedErr: fmt.Errorf("unsupported value type: *protos.valueUnknown"),
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := ConvertFromValue(tc.input)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestConvertToFields(t *testing.T) {
	testCases := []struct {
		input       map[string]any
		expected    []*Field
		expectedErr bool
	}{
		{
			input: map[string]any{
				"key1": "value1",
				"key2": 123,
			},
			expected: []*Field{
				{
					Name:  "key1",
					Value: &Value{Value: &Value_StringValue{StringValue: "value1"}},
				},
				{
					Name:  "key2",
					Value: &Value{Value: &Value_LongValue{LongValue: 123}},
				},
			},
		},
		{
			input: map[string]any{
				"key1": "value1",
				"key2": struct{}{},
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	sortFunc := func(fields []*Field) func(int, int) bool {
		return func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		}
	}

	for _, tc := range testCases {
		result, err := ConvertToFields(tc.input)
		assert.Equal(t, len(tc.expected), len(result))

		sort.Slice(result, sortFunc(result))
		sort.Slice(tc.expected, sortFunc(result))

		for i := range tc.expected {
			assert.EqualExportedValues(t, tc.expected[i], result[i])
		}

		if tc.expectedErr {
			assert.Error(t, err)
		}
	}
}

func TestConvertFromFields(t *testing.T) {
	testCases := []struct {
		input       []*Field
		expected    map[string]any
		expectedErr bool
	}{
		{
			input: []*Field{
				{
					Name:  "key1",
					Value: &Value{Value: &Value_StringValue{StringValue: "value1"}},
				},
				{
					Name:  "key2",
					Value: &Value{Value: &Value_LongValue{LongValue: 123}},
				},
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": int64(123),
			},
		},
		{
			input: []*Field{
				{
					Name:  "key1",
					Value: &Value{Value: &Value_StringValue{StringValue: "value1"}},
				},
				{
					Name:  "key2",
					Value: &Value{Value: &valueUnknown{}},
				},
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		result, err := ConvertFromFields(tc.input)

		assert.Equal(t, tc.expected, result)

		if tc.expectedErr {
			assert.Error(t, err)
		}
	}
}

type unknownVectorType struct{}

func (*unknownVectorType) isVector_Data() {} //nolint:revive,stylecheck // Grpc generated

func TestConvertFromVector(t *testing.T) {
	testCases := []struct {
		input       *Vector
		expected    any
		expectedErr error
	}{
		{
			input: &Vector{
				Data: &Vector_FloatData{
					FloatData: &FloatData{
						Value: []float32{1, 2},
					},
				},
			},
			expected: []float32{float32(1), float32(2)},
		},
		{
			input: &Vector{
				Data: &Vector_BoolData{
					BoolData: &BoolData{
						Value: []bool{true, false},
					},
				},
			},
			expected: []bool{true, false},
		},
		{
			input:       &Vector{Data: &unknownVectorType{}},
			expected:    nil,
			expectedErr: fmt.Errorf("unsupported vector data type: *protos.unknownVectorType"),
		},
	}

	for _, tc := range testCases {
		result, err := ConvertFromVector(tc.input)

		assert.Equal(t, tc.expected, result)
		assert.Equal(t, tc.expectedErr, err)
	}
}
