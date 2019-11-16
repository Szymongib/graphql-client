package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ContentTypeHeader = "Content-Type"
	AcceptHeader      = "Accept"
)

var (
	DefaultParserOptions ParserOptions
)

type Request struct {
	Query  string
	Header http.Header

	// TODO: files
	// TODO: vars
}

func NewRequestRaw(query string, header ...http.Header) Request {
	return Request{
		Query:  query,
		Header: mergeHeaders(header),
	}
}

func NewRequest(operation Operation, headers ...http.Header) (Request, error) {
	return newRequest(operation, DefaultParserOptions, headers...)
}

func newRequest(operation Operation, parserOpts ParserOptions, headers ...http.Header) (Request, error) {
	query, err := operation.ToQueryString(parserOpts)
	if err != nil {
		return Request{}, fmt.Errorf("failed to create GraphQL request from data, %w", err)
	}

	return Request{
		Query:  query,
		Header: mergeHeaders(headers),
	}, nil
}

func (r Request) ToHttpRequest(endpoint string) (*http.Request, error) {
	requestData := r.toRequestData()

	var requestBodyBuffer bytes.Buffer
	if err := json.NewEncoder(&requestBodyBuffer).Encode(requestData); err != nil {
		return nil, fmt.Errorf("failed to encode body: %w", err)
	}

	httpRequest, err := http.NewRequest(http.MethodPost, endpoint, &requestBodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpRequest.Header = r.Header
	r.Header.Set(ContentTypeHeader, "application/json; charset=utf-8")
	r.Header.Set(AcceptHeader, "application/json; charset=utf-8")

	return httpRequest, nil
}

func (r Request) toRequestData() gqlRequestData {
	return gqlRequestData{
		Query:         r.Query,
		OperationName: "",
		Variables:     nil,
	}
}

func mergeHeaders(headers []http.Header) http.Header {
	mergedHeaders := http.Header{}

	for _, header := range headers {
		for h, v := range header {
			mergedHeaders[h] = append(mergedHeaders[h], v...)
		}
	}

	return mergedHeaders
}
