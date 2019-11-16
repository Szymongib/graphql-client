package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleStruct struct {
	StringField    string
	IntField       int
	StringPtrField *string
	BoolField      bool
	InterfaceField interface{}
}

type sliceStruct struct {
	StringsSlice    []string
	IntsSlice       []int
	StringPtrsSlice []*string
	BoolsSlice      []bool
}

type embeddedStruct struct {
	StringField       string
	IntField          int
	SimpleStructField simpleStruct
	SliceStructField  sliceStruct
}

type complexStruct struct {
	StringField          string
	EmbeddedStructField  embeddedStruct
	EmbeddedStructsField []embeddedStruct
	SliceStructsField    []sliceStruct
}

type jsonTaggedStruct struct {
	StringField string   `json:"stringField"`
	IntField    int      `json:"intField"`
	SliceField  []string `json:"sliceField"`
}

type jsonTaggedComplexStruct struct {
	StringField       string             `json:"stringField"`
	IntField          int                `json:"intField"`
	JsonTaggedStruct  jsonTaggedStruct   `json:"jsonTaggedStruct"`
	JsonTaggedStructs []jsonTaggedStruct `json:"jsonTaggedStructs"`
}

type includedPointersStruct struct {
	Name              string          `json:"name"`
	SimpleStructPtr   *simpleStruct   `json:"simpleStructPtr"`
	EmbeddedStructPtr *embeddedStruct `json:"embeddedStructPtr"`
}

type slicesOfPointersStruct struct {
	Name               string            `json:"name"`
	SimpleStructPtrs   []*simpleStruct   `json:"simpleStructPtrs"`
	EmbeddedStructPtrs []*embeddedStruct `json:"embeddedStructPtrs"`
}

type MapAlias map[string]interface{}

type mapsStruct struct {
	SimpleMap         map[string]string
	FloatMap          map[string]float64
	StructMap         map[string]complexStruct
	InterfaceMap      map[string]interface{}
	PointersMap       map[string]*simpleStruct
	MapAlias          MapAlias
	NilMapPlaceholder map[string]string
}

var nilSliceOfSimpleStructs []*simpleStruct

func Test_ParseToGQLQuery(t *testing.T) {

	for _, testCase := range []struct {
		name          string
		data          interface{}
		expectedQuery string
	}{
		{
			name: "simple struct",
			data: simpleStruct{
				StringField:    "test",
				IntField:       1,
				StringPtrField: nil,
				BoolField:      false,
				InterfaceField: "test",
			},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField 
}`,
		},
		{
			name: "simple struct pointer",
			data: &simpleStruct{
				StringField:    "test",
				IntField:       1,
				StringPtrField: nil,
				BoolField:      false,
				InterfaceField: "test",
			},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField 
}`,
		},
		{
			name: "simple struct with struct interface",
			data: simpleStruct{
				StringField:    "test",
				IntField:       1,
				StringPtrField: nil,
				BoolField:      false,
				InterfaceField: simpleStruct{},
			},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
}`,
		},
		{
			name: "simple struct with slice interface",
			data: simpleStruct{
				StringField:    "test",
				IntField:       1,
				StringPtrField: nil,
				BoolField:      false,
				InterfaceField: []simpleStruct{},
			},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
}`,
		},
		{
			name: "slice struct",
			data: sliceStruct{
				StringsSlice:    []string{"nil", ""},
				IntsSlice:       []int{1, 2},
				StringPtrsSlice: []*string{nil},
				BoolsSlice:      []bool{true, false},
			},
			expectedQuery: `{
	StringsSlice 
	IntsSlice 
	StringPtrsSlice 
	BoolsSlice 
}`,
		},
		{
			name: "slice struct with nil slices",
			data: sliceStruct{
				StringsSlice:    nil,
				IntsSlice:       nil,
				StringPtrsSlice: nil,
				BoolsSlice:      nil,
			},
			expectedQuery: `{
	StringsSlice 
	IntsSlice 
	StringPtrsSlice 
	BoolsSlice 
}`,
		},
		{
			name: "embedded struct",
			data: embeddedStruct{
				StringField:       "test",
				IntField:          1,
				SimpleStructField: simpleStruct{},
				SliceStructField:  sliceStruct{},
			},
			expectedQuery: `{
	StringField 
	IntField 
	SimpleStructField {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
	SliceStructField {
		StringsSlice 
		IntsSlice 
		StringPtrsSlice 
		BoolsSlice 
	}
}`,
		},
		{
			name: "embedded struct with filled interfaces",
			data: embeddedStruct{
				StringField: "test",
				IntField:    1,
				SimpleStructField: simpleStruct{
					InterfaceField: simpleStruct{},
				},
				SliceStructField: sliceStruct{},
			},
			expectedQuery: `{
	StringField 
	IntField 
	SimpleStructField {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
	}
	SliceStructField {
		StringsSlice 
		IntsSlice 
		StringPtrsSlice 
		BoolsSlice 
	}
}`,
		},
		{
			name: "complex struct",
			data: complexStruct{
				StringField:          "test",
				EmbeddedStructField:  embeddedStruct{},
				EmbeddedStructsField: []embeddedStruct{},
				SliceStructsField:    []sliceStruct{},
			},
			expectedQuery: `{
	StringField 
	EmbeddedStructField {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
	EmbeddedStructsField {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
	SliceStructsField {
		StringsSlice 
		IntsSlice 
		StringPtrsSlice 
		BoolsSlice 
	}
}`,
		},
		{
			name: "json-tagged struct",
			data: jsonTaggedStruct{
				StringField: "test",
				IntField:    1,
				SliceField:  []string{},
			},
			expectedQuery: `{
	stringField 
	intField 
	sliceField 
}`,
		},
		{
			name: "json-tagged complex struct",
			data: jsonTaggedComplexStruct{
				StringField:       "test",
				IntField:          1,
				JsonTaggedStruct:  jsonTaggedStruct{},
				JsonTaggedStructs: nil,
			},
			expectedQuery: `{
	stringField 
	intField 
	jsonTaggedStruct {
		stringField 
		intField 
		sliceField 
	}
	jsonTaggedStructs {
		stringField 
		intField 
		sliceField 
	}
}`,
		},
		{
			name:          "primitive",
			data:          "abcd",
			expectedQuery: ``,
		},
		{
			name: "slice of pointers",
			data: []*simpleStruct{},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField 
}`,
		},
		{
			name: "nil slice of pointers",
			data: nilSliceOfSimpleStructs,
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField 
}`,
		},
		{
			name: "slice of structs",
			data: []simpleStruct{},
			expectedQuery: `{
	StringField 
	IntField 
	StringPtrField 
	BoolField 
	InterfaceField 
}`,
		},
		{
			name: "included pointers struct with nil pointers",
			data: includedPointersStruct{},
			expectedQuery: `{
	name 
	simpleStructPtr {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
	embeddedStructPtr {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
}`,
		},
		{
			name: "included pointers struct with non nil pointers",
			data: includedPointersStruct{SimpleStructPtr: &simpleStruct{}, EmbeddedStructPtr: &embeddedStruct{}},
			expectedQuery: `{
	name 
	simpleStructPtr {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
	embeddedStructPtr {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
}`,
		},
		{
			name: "slices of pointers struct with non nil slices",
			data: slicesOfPointersStruct{SimpleStructPtrs: []*simpleStruct{}, EmbeddedStructPtrs: []*embeddedStruct{}},
			expectedQuery: `{
	name 
	simpleStructPtrs {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
	embeddedStructPtrs {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
}`,
		},
		{
			name: "slices of pointers struct with nil slices",
			data: slicesOfPointersStruct{SimpleStructPtrs: nil, EmbeddedStructPtrs: nil},
			expectedQuery: `{
	name 
	simpleStructPtrs {
		StringField 
		IntField 
		StringPtrField 
		BoolField 
		InterfaceField 
	}
	embeddedStructPtrs {
		StringField 
		IntField 
		SimpleStructField {
			StringField 
			IntField 
			StringPtrField 
			BoolField 
			InterfaceField 
		}
		SliceStructField {
			StringsSlice 
			IntsSlice 
			StringPtrsSlice 
			BoolsSlice 
		}
	}
}`,
		},
		{
			name: "maps struct",
			data: mapsStruct{},
			expectedQuery: `{
	SimpleMap 
	FloatMap 
	StructMap 
	InterfaceMap 
	PointersMap 
	MapAlias 
	NilMapPlaceholder 
}`,
		},
		{
			name:          "map alias",
			data:          MapAlias{},
			expectedQuery: ``,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			query := ParseToGQLQuery(testCase.data)
			t.Log(query)
			assert.Equal(t, testCase.expectedQuery, query)
		})
	}

}
