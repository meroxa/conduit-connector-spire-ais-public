package ais

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"fmt"
	"net/url"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/machinebox/graphql"
)

type Source struct {
	sdk.UnimplementedSource

	config Config
}

func NewSource() sdk.Source {
	// Create Source and wrap it in the default middleware.
	return sdk.SourceWithMiddleware(&Source{})
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	// Parameters is a map of named Parameters that describe how to configure
	// the Source. Parameters can be generated from SourceConfig with paramgen.
	return map[string]sdk.Parameter{
		"apiUrl": {
			Type:        sdk.ParameterTypeString,
			Default:     "",
			Description: "Url to the Spire AIS GraphQL endpoint",
		},
		"token": {
			Type:        sdk.ParameterTypeString,
			Default:     "",
			Description: "Authentication token for the Spire AIS GraphQL endpoint",
		},
		"query": {
			Type:        sdk.ParameterTypeString,
			Default:     "",
			Description: "GraphQL query string",
		},
	}
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	// Configure is the first function to be called in a connector. It provides
	// the connector with the configuration that can be validated and stored.
	// In case the configuration is not valid it should return an error.
	// Testing if your connector can reach the configured data source should be
	// done in Open, not in Configure.
	// The SDK will validate the configuration and populate default values
	// before calling Configure. If you need to do more complex validations you
	// can do them manually here.

	sdk.Logger(ctx).Info().Msg("Configuring Source...")
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start producing records. If needed, the plugin should open connections in
	// this function. The position parameter will contain the position of the
	// last record that was successfully processed, Source should therefore
	// start producing records after this position. The context passed to Open
	// will be cancelled once the plugin receives a stop signal from Conduit.

	// Create new GraphQL client using URL from config
	graphqlClient := graphql.NewClient(s.config.apiUrl)
	graphqlRequest := graphql.NewRequest(s.config.query)
	// set header fields
	graphqlRequest.Header.Set("Authorization", "Bearer %s")
	err := s.initPosition(pos)
	if err != nil {
		return fmt.Errorf("failed initializing position: %w", err)
	}
	return err

	return graphqlClient
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	// Read returns a new Record and is supposed to block until there is either
	// a new record or the context gets cancelled. It can also return the error
	// ErrBackoffRetry to signal to the SDK it should call Read again with a
	// backoff retry.
	// If Read receives a cancelled context or the context is cancelled while
	// Read is running it must stop retrieving new records from the source
	// system and start returning records that have already been buffered. If
	// there are no buffered records left Read must return the context error to
	// signal a graceful stop. If Read returns ErrBackoffRetry while the context
	// is cancelled it will also signal that there are no records left and Read
	// won't be called again.
	// After Read returns an error the function won't be called again (except if
	// the error is ErrBackoffRetry, as mentioned above).
	// Read can be called concurrently with Ack.
	return sdk.Record{}, nil
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	// Ack signals to the implementation that the record with the supplied
	// position was successfully processed. This method might be called after
	// the context of Read is already cancelled, since there might be
	// outstanding acks that need to be delivered. When Teardown is called it is
	// guaranteed there won't be any more calls to Ack.
	// Ack can be called concurrently with Read.
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that there will be no more calls to any
	// other function. After Teardown returns, the plugin should be ready for a
	// graceful shutdown.
	return nil
}

func (s *Source) validateConfig(cfg map[string]string) error {
	apiUrl, ok := cfg["apiUrl"]
	if !ok {
		return requiredConfigErr("apiUrl")
	}

	_, err := cfg["token"]
	if err {
		return requiredConfigErr("token")
	}

	_, queryErr := cfg["query"]
	if queryErr {
		return requiredConfigErr("query")
	}

	// Check if url is valid
	_, validURLErr = url.ParseRequestURI(apiUrl)
	if validURLErr {
		return fmt.Errorf("%q is not a valid URL", apiUrl)
	}

	// // make sure we can stat the file, we don't care if it doesn't exist though
	// _, err := os.Stat(path)
	// if err != nil && !os.IsNotExist(err) {
	// 	return fmt.Errorf(
	// 		"%q config value %q does not contain a valid path: %w",
	// 		ConfigPath, path, err,
	// 	)
	// }

	return nil
}
