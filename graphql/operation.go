package graphql

import (
	"fmt"
	"sort"
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

type OperationInput map[string]interface{}

func (o OperationInput) sortedKeys() []string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
