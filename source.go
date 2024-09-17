// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ais

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"fmt"

	"github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/lang"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/machinebox/graphql"
)

type IteratorCreator interface {
	NewIterator(client GraphQLClient, token string, query string, batchSize int, p opencdc.Position) (*Iterator, error)
}

type SourceIteratorCreator struct {
}

func (ic SourceIteratorCreator) NewIterator(client GraphQLClient, token string, query string, batchSize int, p opencdc.Position) (*Iterator, error) {
	return NewIterator(client, token, query, batchSize, p)
}

type Source struct {
	sdk.UnimplementedSource

	config               SourceConfig
	iterator             *Iterator
	iteratorCreator      IteratorCreator
	startQueryFromCursor bool
}

type SourceConfig struct {
	// Config includes parameters that are the same in the source and destination.
	Config

	// Query is the GraphQL Query to use when pulling data from the Spire API.
	Query string `json:"query"`
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{
		config:          SourceConfig{},
		iteratorCreator: SourceIteratorCreator{},
	}, sdk.DefaultSourceMiddleware( // disable schema extraction by default, because the source produces raw data
		sdk.SourceWithSchemaExtractionConfig{
			PayloadEnabled: lang.Ptr(false),
			KeyEnabled:     lang.Ptr(false),
		})...)
}

func (s *Source) Parameters() config.Parameters {
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg config.Config) error {
	sdk.Logger(ctx).Debug().Msg("Configuring Source connector...")
	err := sdk.Util.ParseConfig(ctx, cfg, &s.config, NewSource().Parameters())
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if s.config.Query == "" {
		s.config.Query = vesselQuery()
	}

	return nil
}

func (s *Source) Open(ctx context.Context, pos opencdc.Position) error {
	sdk.Logger(ctx).Debug().Msg("Opening Source connector...")
	c := graphql.NewClient(s.config.APIURL)
	it, err := s.iteratorCreator.NewIterator(c, s.config.Token, s.config.Query, s.config.BatchSize, pos)
	if err != nil {
		return fmt.Errorf("failed to create iterator: %w", err)
	}
	s.iterator = it

	if s.iterator.position != nil {
		s.startQueryFromCursor = true
	}
	return nil
}

func (s *Source) Read(ctx context.Context) (opencdc.Record, error) {
	if !s.iterator.HasNext(ctx) && s.iterator.position != nil && !s.startQueryFromCursor {
		return opencdc.Record{}, sdk.ErrBackoffRetry
	}

	record, err := s.iterator.Next(ctx)
	sdk.Logger(context.Background()).Info().Msgf("Nodes processed: %d", s.iterator.nodesProcessed)
	if err != nil {
		return opencdc.Record{}, fmt.Errorf("error reading next record: %w", err)
	}
	return record, nil
}

func (s *Source) Ack(ctx context.Context, position opencdc.Position) error {
	// Ack signals to the implementation that the record with the supplied
	// position was successfully processed. This method might be called after
	// the context of Read is already cancelled, since there might be
	// outstanding acks that need to be delivered. When Teardown is called it is
	// guaranteed there won't be any more calls to Ack.
	// Ack can be called concurrently with Read.
	// sdk.Logger(ctx).Debug().Str("position", string(position)).Msg("got ack")
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that there will be no more calls to any
	// other function. After Teardown returns, the plugin should be ready for a
	// graceful shutdown.

	return nil
}

func (s *Source) GetConfig() SourceConfig {
	return s.config
}
