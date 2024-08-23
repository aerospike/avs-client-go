package protos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetStrPtr(str string) *string {
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
				Set:       GetStrPtr("testSet"),
				Value: &Key_StringValue{
					StringValue: "testString",
				},
			},
		},
		{
			input: []byte{0x01, 0x02, 0x03},
			expected: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value: &Key_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
		},
		{
			input: int32(123),
			expected: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value: &Key_IntValue{
					IntValue: 123,
				},
			},
		},
		{
			input: int64(123456789),
			expected: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
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
		result, err := ConvertToKey("testNamespace", GetStrPtr("testSet"), tc.input)

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
				Set:       GetStrPtr("testSet"),
				Value: &Key_StringValue{
					StringValue: "testString",
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       GetStrPtr("testSet"),
			expectedKey:       "testString",
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value: &Key_BytesValue{
					BytesValue: []byte{0x01, 0x02, 0x03},
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       GetStrPtr("testSet"),
			expectedKey:       []byte{0x01, 0x02, 0x03},
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value: &Key_IntValue{
					IntValue: 123,
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       GetStrPtr("testSet"),
			expectedKey:       int32(123),
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value: &Key_LongValue{
					LongValue: 123456789,
				},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       GetStrPtr("testSet"),
			expectedKey:       int64(123456789),
		},
		{
			input: &Key{
				Namespace: "testNamespace",
				Set:       GetStrPtr("testSet"),
				Value:     &keyUnknown{},
			},
			expectedNamespace: "testNamespace",
			expectedSet:       GetStrPtr("testSet"),
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
		expectedErr error
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
			input:       struct{}{}, // Unsupported type
			expected:    nil,
			expectedErr: fmt.Errorf("unsupported value type: struct {}"),
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result, err := ConvertToValue(tc.input)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedErr, err)
		})
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
