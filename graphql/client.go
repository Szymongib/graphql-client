package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type gqlRequestData struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type gqlResponseData struct {
	Data   interface{} `json:"data"`
	Errors []graphErr  `json:"errors"`
}

type graphErr struct {
	Message string
}

type resultWrapper struct {
	Result interface{} `json:"result"`
}

func (e graphErr) Error() string {
	return "graphql error: " + e.Message
}

type Client struct {
	*options
	endpoint string
}

func NewClient(endpoint string, option ...Option) *Client {
	options := &options{
		parserOptions: DefaultParserOptions,
		httpClient:    &http.Client{},
		logger:        nil,
	}

	for _, opt := range option {
		opt.apply(options)
	}

	return &Client{
		endpoint: endpoint,
		options:  options,
	}
}

func (c Client) Execute(ctx context.Context, request Request, responseOut interface{}) error {
	err := checkContext(ctx)
	if err != nil {
		return err
	}

	return c.executeRequest(ctx, request, responseOut)
}

// TODO - try to get rid of result here?

// Run executes GraphQL operation
func (c Client) Run(ctx context.Context, operation Operation, result interface{}, header ...http.Header) error {
	if operation.Requested == nil {
		operation.Requested = result
	}

	request, err := newRequest(operation, c.options.parserOptions, header...)
	if err != nil {
		return err
	}

	return c.wrapAndExecute(ctx, request, result)
}

// Query executes GraphQL query parsing provided input and requested type to query string
func (c Client) Query(ctx context.Context, name string, input OperationInput, requested interface{}, header ...http.Header) error {
	operation := Operation{
		Type:      Query,
		Name:      name,
		Requested: requested,
		Input:     input,
	}

	return c.Run(ctx, operation, &requested, header...)
}

// Mutate executes GraphQL mutation parsing provided input and requested type to query string
func (c Client) Mutate(ctx context.Context, name string, input OperationInput, requested interface{}, header ...http.Header) error {
	operation := Operation{
		Type:      Mutation,
		Name:      name,
		Requested: requested,
		Input:     input,
	}

	return c.Run(ctx, operation, &requested, header...)
}

func (c Client) wrapAndExecute(ctx context.Context, request Request, result interface{}) error {
	resultWrapper := resultWrapper{Result: result}
	return c.executeRequest(ctx, request, &resultWrapper)
}

func (c Client) executeRequest(ctx context.Context, request Request, responseOut interface{}) error {
	c.logRequest(request)

	httpRequest, err := request.ToHttpRequest(c.endpoint)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	httpRequest = httpRequest.WithContext(ctx)

	res, err := c.options.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("error while executing request: %w", err)
	}
	defer c.closeResponse(res.Body)

	c.logResponse(res)

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil || len(bodyBytes) == 0 {
			return fmt.Errorf("received unexpected response status: %s", res.Status)
		}

		errorResponse := gqlResponseData{}
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil || len(errorResponse.Errors) == 0 {
			return fmt.Errorf("received unexpected response status: %s. Response body: %s", res.Status, string(bodyBytes))
		}

		return parseErrorResponse(errorResponse.Errors)
	}

	responseData := gqlResponseData{
		Data:   responseOut,
		Errors: []graphErr{},
	}

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(responseData.Errors) > 0 {
		return parseErrorResponse(responseData.Errors)
	}

	return nil
}

func parseErrorResponse(gqlErrors []graphErr) error {
	errorString := ""

	for _, err := range gqlErrors {
		errorString += fmt.Sprintf("%s%s\n", errorString, err.Error())
	}

	return fmt.Errorf(strings.TrimSuffix(errorString, "\n"))
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (c Client) logRequest(request Request) {
	if c.options.logger != nil {
		c.options.logger(fmt.Sprintf("Executing request\nHeaders: %s\nQuery: %s", request.Header, request.Query))
	}
}

func (c Client) logResponse(response *http.Response) {
	if c.options.logger != nil {
		r1, r2, err := drainBody(response.Body)
		if err != nil {
			c.options.logger("Failed to log response body: failed to copy response body")
			return
		}
		response.Body = r1

		buff := make([]byte, response.ContentLength)
		_, err = r2.Read(buff)
		if err != nil {
			c.options.logger("Failed to log response body: failed to read response body copy")
			return
		}

		c.options.logger(fmt.Sprintf("Response status: %s\nBody:\n%s", response.Status, string(buff)))
	}
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func (c Client) closeResponse(closer io.ReadCloser) {
	err := closer.Close()
	if err != nil {
		c.options.logger(fmt.Sprintf("Warning: failed to close response body: %s", err.Error()))
	}
}
