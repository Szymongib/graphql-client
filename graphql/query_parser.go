package graphql

import (
	"fmt"
	"reflect"
)

// TODO: try to come up with a mechanism to handle unions
// TODO: support some tags to allow user to skip some fields

const (
	jsonTagKey = "json"
)

func ParseToGQLQuery(data interface{}, nestedInputs ...NestedOperationInput) string {
	return parseToGQLQuery(data, "", 0, nestedInputs)
}

func parseToGQLQuery(data interface{}, fieldPath FieldPath, indent int, nestedInputs []NestedOperationInput) string {
	reflectVal := reflect.ValueOf(data)

	reflectVal = unwrapPointerOrInterface(reflectVal)

	if reflectVal.Kind() == reflect.Struct {
		fieldsString := "{"

		for i := 0; i < reflectVal.NumField(); i++ {
			queriedName := reflectVal.Type().Field(i).Tag.Get(jsonTagKey)
			if queriedName == "" {
				queriedName = reflectVal.Type().Field(i).Name
			}

			inputString := ""
			currentFieldPath := fieldPath.Append(queriedName)
			for _, input := range nestedInputs {
				if currentFieldPath.Matches(input.FieldPath) {
					var err error
					inputString, err = ParseToGQLInput(input.Input)
					if err != nil {
						panic(err) // TODO - handle error
					}

					inputString = fmt.Sprintf("(%s)", inputString)
				}
			}

			fieldQuery := parseToGQLQuery(reflectVal.Field(i).Interface(), currentFieldPath, indent+1, nestedInputs)

			field := queriedName
			field = appendIfNotEmpty(field, inputString, "")
			field = appendIfNotEmpty(field, fieldQuery, " ")

			fieldsString = fmt.Sprintf("%s\n%s%s", fieldsString, tabsIndent(indent+1), field)
		}

		fieldsString += "\n" + tabsIndent(indent) + "}"

		return fieldsString
	}

	if reflectVal.Kind() == reflect.Slice || reflectVal.Kind() == reflect.Array {
		sliceElemObj := reflectVal.Type().Elem()
		return parseToGQLQuery(reflect.New(sliceElemObj).Interface(), fieldPath, indent, nestedInputs)
	}

	return ""
}

func unwrapPointerOrInterface(reflectVal reflect.Value) reflect.Value {
	for reflectVal.Kind() == reflect.Ptr || reflectVal.Kind() == reflect.Interface {
		reflectValElem := reflectVal.Elem()

		// In case that the kind was Ptr and the value is nil then the Elem().Kind() is invalid
		// to inspect the fields we need to instantiate new object
		if reflectValElem.Kind() == reflect.Invalid {
			reflectVal = reflect.New(reflectVal.Type().Elem())
		} else {
			reflectVal = reflectValElem
		}
	}

	return reflectVal
}

func tabsIndent(tabsCount int) (tabs string) {
	for i := 0; i < tabsCount; i++ {
		tabs += "\t"
	}

	return tabs
}

func appendIfNotEmpty(base, appendix, separator string) string {
	if appendix != "" {
		return fmt.Sprintf("%s%s%s", base, separator, appendix)
	}
	return base
}
