package graphql

import (
	"fmt"
	"sort"
	"strings"
)

type Operation struct {
	Type      OperationType
	Name      string
	Requested interface{}
	Input     OperationInput
}

func (o Operation) ToQueryString(options ...ParserOptions) (string, error) {
	parsedInput, err := ParseToGQLInput(o.Input, options...)
	if err != nil {
		return "", fmt.Errorf("failed to create query string, %w", err)
	}

	return fmt.Sprintf(`%s {
	result: %s(%s) %s
}`, o.Type, o.Name, parsedInput, ParseToGQLQuery(o.Requested)), nil
}

type OperationType string

const (
	Query    OperationType = "query"
	Mutation OperationType = "mutation"
)

type Input struct {
	// Input is input of main operation
	Input OperationInput
	// NestedInputs are inputs of sub operations
	NestedInputs []NestedOperationInput
}

type OperationInput map[string]interface{}

type NestedOperationInput struct {
	// FieldPath is jsonpath of the field to which query input should be used
	FieldPath FieldPath
	// Input is string containing GraphQL query input in from of comma separated parameters, ex:
	// `myParam: "paramValue", otherParam: { first: "firstValue" second: "secondValue"}`
	// can be created by ParseToGraphQLInput
	Input string
}

type FieldPath string

func (fp FieldPath) Append(name string) FieldPath {
	return FieldPath(fmt.Sprintf("%s.%s", fp, name))
}

func (fp FieldPath) Matches(others ...FieldPath) bool {
	trimmed := strings.TrimPrefix(string(fp), ".")

	for _, otherFP := range others {
		otherTrimmed := strings.TrimPrefix(string(otherFP), ".")

		if trimmed == otherTrimmed {
			return true
		}
	}

	return false
}

func (o OperationInput) sortedKeys() []string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
