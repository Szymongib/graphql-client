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

func ParseToGQLQuery(data interface{}) string {
	return parseToGQLQuery(data, 0)
}

func parseToGQLQuery(data interface{}, indent int) string {
	reflectVal := reflect.ValueOf(data)

	reflectVal = unwrapPointerOrInterface(reflectVal)

	if reflectVal.Kind() == reflect.Struct {
		fieldsString := "{"

		for i := 0; i < reflectVal.NumField(); i++ {
			queriedName := reflectVal.Type().Field(i).Tag.Get(jsonTagKey)
			if queriedName == "" {
				queriedName = reflectVal.Type().Field(i).Name
			}

			// TODO: this space may be confusing, consider removing it
			field := fmt.Sprintf("%s %s", queriedName, parseToGQLQuery(reflectVal.Field(i).Interface(), indent+1))
			fieldsString = fmt.Sprintf("%s\n%s%s", fieldsString, tabsIndent(indent+1), field)
		}

		fieldsString += "\n" + tabsIndent(indent) + "}"

		return fieldsString
	}

	if reflectVal.Kind() == reflect.Slice || reflectVal.Kind() == reflect.Array {
		sliceElemObj := reflectVal.Type().Elem()
		return parseToGQLQuery(reflect.New(sliceElemObj).Interface(), indent)
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
