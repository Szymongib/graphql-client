package graphql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: consider making it as a parser struct not as a function - can have both query parser and input parser

type ParserOptions struct {
	// SkipZeroValues determines if parameters with zero values should be skipped when parsing to input
	// Be aware that values like 'false' for the bool field or "" for a string field are also zero values
	SkipZeroValues bool
}

func ParseToGQLInput(input OperationInput, options ...ParserOptions) (string, error) {
	if input == nil {
		return "", nil
	}

	var opts = ParserOptions{}

	if len(options) != 0 {
		opts = options[0]
	}

	return opts.parseToGQLInput(input)
}

func (o ParserOptions) parseToGQLInput(input OperationInput) (string, error) {
	gqlInput := ""
	sortedKeys := input.sortedKeys()

	for _, paramName := range sortedKeys {
		value := input[paramName]

		inputValue, ok, err := o.objectToGQLInput(value, 0)
		if err != nil {
			return "", fmt.Errorf("failed to parse to GQL input, %w", err)
		}
		if !ok {
			return "", fmt.Errorf("failed to parse to GQL input, invalid input for %s parameter", paramName)
		}
		gqlInput = fmt.Sprintf("%s %s: %s,", gqlInput, paramName, inputValue)
	}

	gqlInput = strings.TrimPrefix(gqlInput, " ")
	return strings.TrimSuffix(gqlInput, ","), nil
}

func (o ParserOptions) objectToGQLInput(object interface{}, indent int) (string, bool, error) {
	reflectVal := reflect.ValueOf(object)

	if reflectVal.Kind() == reflect.Invalid {
		return "", false, nil
	}

	if o.SkipZeroValues && reflectVal.IsZero() {
		return "", false, nil
	}

	switch reflectVal.Kind() {
	case reflect.Struct:
		return o.structToGQLInput(reflectVal, indent)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := strconv.FormatInt(reflectVal.Int(), 10)
		return value, true, nil
	case reflect.String:
		return o.stringToGQLInput(reflectVal), true, nil
	case reflect.Float32:
		value := strconv.FormatFloat(reflectVal.Float(), 'f', 8, 32)
		return value, true, nil
	case reflect.Float64:
		value := strconv.FormatFloat(reflectVal.Float(), 'f', 20, 64)
		return value, true, nil
	case reflect.Bool:
		return strconv.FormatBool(reflectVal.Bool()), true, nil
	case reflect.Ptr, reflect.Interface:
		reflectValElem := reflectVal.Elem()
		if !reflectValElem.IsValid() {
			return "", false, nil
		}
		return o.objectToGQLInput(reflectValElem.Interface(), indent)
	case reflect.Array, reflect.Slice:
		return o.arrayToGQLInput(reflectVal, indent)
	case reflect.Map:
		return o.mapToGQLInput(reflectVal, indent)
	}

	// TODO - other cases (if need to support more?)

	return "", false, nil
}

func (o ParserOptions) stringToGQLInput(reflectVal reflect.Value) string {
	return fmt.Sprintf("\"%s\"", reflectVal.String())
}

func (o ParserOptions) structToGQLInput(reflectVal reflect.Value, indent int) (string, bool, error) {
	fieldsString := ""

	for i := 0; i < reflectVal.NumField(); i++ {
		inputName := reflectVal.Type().Field(i).Tag.Get(jsonTagKey)
		if inputName == "" {
			inputName = reflectVal.Type().Field(i).Name
		}

		inputValue, ok, err := o.objectToGQLInput(reflectVal.Field(i).Interface(), indent+1)
		if err != nil {
			return "", false, err
		}
		if ok {
			fieldsString = fmt.Sprintf("%s\n%s%s: %s", fieldsString, tabsIndent(indent+1), inputName, inputValue)
		}
	}

	if fieldsString == "" {
		return "", false, nil
	}

	fieldsString = "{" + fieldsString + "\n" + tabsIndent(indent) + "}"

	return fieldsString, true, nil
}

func (o ParserOptions) arrayToGQLInput(reflectVal reflect.Value, indent int) (string, bool, error) {
	if reflectVal.IsNil() {
		return "", false, nil
	}

	elementsString := "["
	for i := 0; i < reflectVal.Len(); i++ {
		arrayElem := reflectVal.Index(i)

		inputValue, ok, err := o.objectToGQLInput(arrayElem.Interface(), indent+1)
		if err != nil {
			return "", false, err
		}
		if ok {
			elem := fmt.Sprintf("\t%s", inputValue)
			elementsString = fmt.Sprintf("%s\n%s%s,", elementsString, tabsIndent(indent), elem)
		}
	}

	elementsString = strings.TrimSuffix(elementsString, ",")
	elementsString += "\n" + tabsIndent(indent) + "]"

	return elementsString, true, nil
}

func (o ParserOptions) mapToGQLInput(reflectVal reflect.Value, indent int) (string, bool, error) {
	if reflectVal.IsNil() {
		return "", false, nil
	}

	mapElemsString := "{"

	mapIter := reflectVal.MapRange()
	for mapIter.Next() {
		key := mapIter.Key()
		if key.Kind() != reflect.String {
			return "", false, fmt.Errorf("unsupported map key type %s, must be of kind string", key.Kind())
		}

		value, ok, err := o.objectToGQLInput(mapIter.Value().Interface(), indent+1)
		if err != nil {
			return "", false, err
		} else if !ok {
			return "", false, nil
		}

		// TODO - make sure that key should not be in ""
		elem := fmt.Sprintf("\t%s: %s", key, value)
		mapElemsString = fmt.Sprintf("%s\n%s%s", mapElemsString, tabsIndent(indent), elem)
	}
	mapElemsString += "\n" + tabsIndent(indent) + "}"

	return mapElemsString, true, nil
}
