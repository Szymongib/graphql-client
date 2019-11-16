package graphql

import (
	"fmt"
	"testing"

	"github.com/szymongib/graphql-client/util"

	"github.com/stretchr/testify/assert"
)

func TestParseToGQLInput(t *testing.T) {

	for _, testCase := range []struct {
		description   string
		input         OperationInput
		expectedInput string
	}{
		{
			description: "simple struct input with nils",
			input: OperationInput{
				"in": simpleStruct{
					StringField:    "test",
					IntField:       3,
					StringPtrField: nil,
					BoolField:      false,
					InterfaceField: nil,
				},
			},
			expectedInput: `in: {
	StringField: "test"
	IntField: 3
	BoolField: false
}`,
		},
		{
			description: "simple struct input with non nils",
			input: OperationInput{
				"in": simpleStruct{
					StringField:    "test",
					IntField:       3,
					StringPtrField: util.StringPtr("test-ptr"),
					BoolField:      false,
					InterfaceField: simpleStruct{StringField: "inner-test"},
				},
			},
			expectedInput: `in: {
	StringField: "test"
	IntField: 3
	StringPtrField: "test-ptr"
	BoolField: false
	InterfaceField: {
		StringField: "inner-test"
		IntField: 0
		BoolField: false
	}
}`,
		},
		{
			description: "simple struct input with non nils",
			input: OperationInput{
				"in": simpleStruct{
					StringField:    "test",
					IntField:       3,
					StringPtrField: util.StringPtr("test-ptr"),
					BoolField:      false,
					InterfaceField: simpleStruct{StringField: "inner-test"},
				},
			},
			expectedInput: `in: {
	StringField: "test"
	IntField: 3
	StringPtrField: "test-ptr"
	BoolField: false
	InterfaceField: {
		StringField: "inner-test"
		IntField: 0
		BoolField: false
	}
}`,
		},
		{
			description: "slice of pointers struct",
			input: OperationInput{
				"in": slicesOfPointersStruct{
					Name: "test",
					SimpleStructPtrs: []*simpleStruct{
						{
							StringField:    "first",
							IntField:       1,
							StringPtrField: nil,
							BoolField:      true,
							InterfaceField: "test",
						},
						{
							StringField:    "second",
							IntField:       2,
							StringPtrField: util.StringPtr("test"),
							BoolField:      true,
							InterfaceField: nil,
						},
					},
					EmbeddedStructPtrs: nil,
				},
			},
			expectedInput: `in: {
	name: "test"
	simpleStructPtrs: [
		{
			StringField: "first"
			IntField: 1
			BoolField: true
			InterfaceField: "test"
		},
		{
			StringField: "second"
			IntField: 2
			StringPtrField: "test"
			BoolField: true
		}
	]
}`,
		},
		{
			description: "multiple input params",
			input: OperationInput{
				"in": embeddedStruct{
					StringField:       "",
					IntField:          0,
					SimpleStructField: simpleStruct{},
					SliceStructField: sliceStruct{
						StringsSlice:    []string{"test", "test2"},
						IntsSlice:       []int{1, 2},
						StringPtrsSlice: []*string{util.StringPtr("test")},
						BoolsSlice:      []bool{true, true, false},
					},
				},
				"id":    "abcd-efgh",
				"limit": 20,
			},
			expectedInput: `id: "abcd-efgh", in: {
	StringField: ""
	IntField: 0
	SimpleStructField: {
		StringField: ""
		IntField: 0
		BoolField: false
	}
	SliceStructField: {
		StringsSlice: [
			"test",
			"test2"
		]
		IntsSlice: [
			1,
			2
		]
		StringPtrsSlice: [
			"test"
		]
		BoolsSlice: [
			true,
			true,
			false
		]
	}
}, limit: 20`,
		},
		{
			description: "arrays as a params",
			input: OperationInput{
				"floats":  []float32{21.2, 64.3534},
				"strings": []string{"aaa", "bbb"},
				"structs": []simpleStruct{{
					StringField:    "test",
					IntField:       0,
					StringPtrField: util.StringPtr("test-ptr"),
					BoolField:      false,
				}},
				"pointers": []*simpleStruct{{
					StringField:    "test",
					IntField:       0,
					StringPtrField: util.StringPtr("test-ptr"),
					BoolField:      false,
				}, nil, nil},
			},
			expectedInput: `floats: [
	21.20000076,
	64.35340118
], pointers: [
	{
		StringField: "test"
		IntField: 0
		StringPtrField: "test-ptr"
		BoolField: false
	}
], strings: [
	"aaa",
	"bbb"
], structs: [
	{
		StringField: "test"
		IntField: 0
		StringPtrField: "test-ptr"
		BoolField: false
	}
]`,
		},
		{
			description: "maps in params",
			input: OperationInput{
				"in": mapsStruct{
					SimpleMap: map[string]string{"k1": "v1"},
					FloatMap:  map[string]float64{"k1": 20.000203},
					StructMap: map[string]complexStruct{"k1": {
						StringField:          "test",
						EmbeddedStructField:  embeddedStruct{},
						EmbeddedStructsField: nil,
						SliceStructsField:    nil,
					}},
					InterfaceMap:      map[string]interface{}{"k4": map[string]interface{}{"kk2": simpleStruct{}}},
					PointersMap:       map[string]*simpleStruct{"k2": nil},
					MapAlias:          map[string]interface{}{"k1": "test"},
					NilMapPlaceholder: nil,
				},
			},
			expectedInput: `in: {
	SimpleMap: {
		k1: "v1"
	}
	FloatMap: {
		k1: 20.00020299999999906504
	}
	StructMap: {
		k1: {
			StringField: "test"
			EmbeddedStructField: {
				StringField: ""
				IntField: 0
				SimpleStructField: {
					StringField: ""
					IntField: 0
					BoolField: false
				}
			}
		}
	}
	InterfaceMap: {
		k4: {
			kk2: {
				StringField: ""
				IntField: 0
				BoolField: false
			}
		}
	}
	MapAlias: {
		k1: "test"
	}
}`,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {

			gqlInput, err := ParseToGQLInput(testCase.input)
			assert.NoError(t, err)

			fmt.Println(gqlInput)

			assert.Equal(t, testCase.expectedInput, gqlInput)

		})
	}

	for _, testCase := range []struct {
		description   string
		input         OperationInput
		expectedInput string
	}{
		{
			description: "simple struct input with nils",
			input: OperationInput{
				"in": simpleStruct{
					StringField:    "test",
					IntField:       3,
					StringPtrField: nil,
					BoolField:      false,
					InterfaceField: nil,
				},
			},
			expectedInput: `in: {
	StringField: "test"
	IntField: 3
}`,
		},
		{
			description: "simple struct with zero values",
			input: OperationInput{
				"in": embeddedStruct{
					StringField:       "test",
					IntField:          0,
					SimpleStructField: simpleStruct{},
					SliceStructField:  sliceStruct{},
				},
			},
			expectedInput: `in: {
	StringField: "test"
}`,
		},
	} {
		t.Run("skip zero values "+testCase.description, func(t *testing.T) {
			gqlInput, err := ParseToGQLInput(testCase.input, ParserOptions{SkipZeroValues: true})
			assert.NoError(t, err)

			fmt.Println(gqlInput)

			assert.Equal(t, testCase.expectedInput, gqlInput)

		})
	}

	t.Run("should return error on nil input", func(t *testing.T) {
		gqlInput, err := ParseToGQLInput(OperationInput{"test": "test", "object": nil})
		assert.Error(t, err)
		assert.Empty(t, gqlInput)
	})

	t.Run("should return error when parsed input is empty", func(t *testing.T) {
		gqlInput, err := ParseToGQLInput(OperationInput{"in": simpleStruct{}}, ParserOptions{SkipZeroValues: true})
		assert.Error(t, err)
		assert.Empty(t, gqlInput)
	})

	t.Run("should return error if input contains unsupported type", func(t *testing.T) {
		unsupportedInput := struct {
			Name string
			Map  map[simpleStruct]int
		}{
			Name: "test",
			Map:  map[simpleStruct]int{simpleStruct{StringField: "test"}: 10},
		}

		gqlInput, err := ParseToGQLInput(OperationInput{"in": unsupportedInput})
		assert.Error(t, err)
		assert.Empty(t, gqlInput)
	})

}
