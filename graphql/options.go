package graphql

import "net/http"

type logger func(string)

type options struct {
	parserOptions ParserOptions
	httpClient    *http.Client
	logger        logger
}

// Option overrides behavior of GraphQLClient.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithHTTPClient(client *http.Client) Option {
	return optionFunc(func(o *options) {
		o.httpClient = client
	})
}

func WithLogger(logger func(string)) Option {
	return optionFunc(func(o *options) {
		o.logger = logger
	})
}

func WithParserOptions(parserOpts ParserOptions) Option {
	return optionFunc(func(o *options) {
		o.parserOptions = parserOpts
	})
}

// TODO: some hook on response copy?
