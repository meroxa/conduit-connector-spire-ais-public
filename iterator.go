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

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/machinebox/graphql"
)

type IteratorInterface interface {
	HasNext(ctx context.Context) bool
	Next(ctx context.Context) (sdk.Record, error)
}

// Add GraphQLClient interface for dependency injection
type GraphQLClient interface {
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

// Updated Iterator struct with logger and client dependencies
type Iterator struct {
	query        string
	token        string
	batchSize    int
	cursor       string
	hasNext      bool
	position     []byte
	client       GraphQLClient
	currentBatch []Node
}

func NewIterator(client GraphQLClient, token string, query string, batchSize int, p sdk.Position) (*Iterator, error) {
	return &Iterator{
		token:     token,
		query:     query,
		batchSize: batchSize,
		client:    client,
		position:  p,
	}, nil
}

// Ensure Iterator implements IteratorInterface
var _ IteratorInterface = (*Iterator)(nil)

func (it *Iterator) HasNext(ctx context.Context) bool {
	// return early if there are more nodes
	if len(it.currentBatch) > 0 {
		return true
	}

	if it.hasNext {
		err := it.loadBatch(ctx)
		if err != nil {
			sdk.Logger(ctx).Err(err)
			return false
		}
		return true
	}

	return false
}

func (it *Iterator) Next(ctx context.Context) (sdk.Record, error) {
	// return next message from cached batch
	var out Node
	if len(it.currentBatch) > 0 {
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	} else {
		err := it.loadBatch(ctx)
		if err != nil {
			sdk.Logger(ctx).Err(err).Msg("loadBatch returned error")
			return sdk.Record{}, fmt.Errorf("loadBatch returned error: %w", err)
		}
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	}

	return wrapAsRecord(out, it.position)
}

// Updated loadBatch function with dependency injection
func (it *Iterator) loadBatch(ctx context.Context) error {
	graphqlRequest := graphql.NewRequest(it.query)
	graphqlRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))
	graphqlRequest.Var("first", it.batchSize)
	var Response struct {
		Vessels Vessels
	}

	lastSuccessfulCursor := it.cursor

	if it.hasNext {
		graphqlRequest.Var("after", it.cursor)
	}

	maxRetries := 3
	retryDelay := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		err := it.client.Run(ctx, graphqlRequest, &Response)
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			sdk.Logger(context.Background()).Err(err).Msg("Retrying query...")
			time.Sleep(retryDelay)
		} else {
			sdk.Logger(context.Background()).Err(err).Msg("%w")
			it.cursor = lastSuccessfulCursor
			return fmt.Errorf("error making graphQL Request: %w", err)
		}
	}

	it.currentBatch = Response.Vessels.Nodes
	it.hasNext = Response.Vessels.PageInfo.HasNextPage
	it.cursor = Response.Vessels.PageInfo.EndCursor
	it.position = []byte(Response.Vessels.PageInfo.EndCursor)

	return nil
}

func wrapAsRecord(in Node, endCursor sdk.Position) (sdk.Record, error) {
	updateTimestamp, err := time.Parse(time.RFC3339, in.UpdateTimestamp)
	if err != nil {
		sdk.Logger(context.Background()).Err(err).Msg("%w")
		return sdk.Record{}, fmt.Errorf("error occurred while wrapping results as a Conduit Record: %w", err)
	}

	sdkMetadata := make(sdk.Metadata)
	sdkMetadata.SetCreatedAt(updateTimestamp)

	b, err := json.Marshal(in)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error occurred marshalling JSON: %w", err)
	}

	return sdk.Util.Source.NewRecordCreate(endCursor, sdkMetadata, nil, sdk.RawData(b)), nil
}
