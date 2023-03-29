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
	"errors"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/machinebox/graphql"
	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
)

type MockGraphQLClient struct {
	RunFn func(ctx context.Context, req *graphql.Request, resp interface{}) error
}

func (m *MockGraphQLClient) Run(ctx context.Context, req *graphql.Request, resp interface{}) error {
	return m.RunFn(ctx, req, resp)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func TestIterator(t *testing.T) {
	t.Run("NewIterator", func(t *testing.T) {
		is := is.New(t)
		client := &graphql.Client{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)

		is.NoErr(err)
		is.Equal(client, it.client)
		is.Equal(token, it.token)
		is.Equal(query, it.query)
		// is.Equal(position, it.position)
	})

	t.Run("HasNext", func(t *testing.T) {
		is := is.New(t)
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)
		is.NoErr(err)

		// Set up expected behavior
		expectedResponse := struct {
			Vessels Vessels
		}{
			Vessels: Vessels{
				PageInfo: PageInfo{HasNextPage: false, EndCursor: ""},
				Nodes:    []Node{},
			},
		}
		client.RunFn = func(ctx context.Context, req *graphql.Request, resp interface{}) error {
			arg := resp.(*struct{ Vessels Vessels })
			*arg = expectedResponse
			return nil
		}

		// hasNext is true
		it.hasNext = true
		is.Equal(true, it.HasNext(context.Background()))

		// hasNext is false, but there are more nodes
		it.hasNext = false
		it.currentBatch = []Node{{}, {}}
		is.Equal(true, it.HasNext(context.Background()))

		// hasNext is false, and no more nodes
		it.currentBatch = []Node{}
		is.Equal(false, it.HasNext(context.Background()))
	})

	t.Run("Next", func(t *testing.T) {
		is := is.New(t)
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)
		is.NoErr(err)

		it.currentBatch = []Node{
			{
				UpdateTimestamp: "2021-10-01T15:00:00Z",
			},
		}

		record, err := it.Next(context.Background())
		is.NoErr(err)
		is.True(record.Payload.After != nil)
	})

	t.Run("loadBatch_HappyPath", func(t *testing.T) {
		is := is.New(t)
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)
		is.NoErr(err)

		// Mock the GraphQL response
		mockResponse := struct {
			Vessels Vessels
		}{
			Vessels: Vessels{
				Nodes: []Node{},
				PageInfo: PageInfo{
					HasNextPage: false,
					EndCursor:   "some_cursor",
				},
			},
		}

		client.RunFn = func(ctx context.Context, req *graphql.Request, resp interface{}) error {
			arg := resp.(*struct{ Vessels Vessels })
			*arg = mockResponse
			return nil
		}

		err = it.loadBatch(context.Background())

		is.NoErr(err)
		is.Equal(mockResponse.Vessels.Nodes, it.currentBatch)
		is.Equal(mockResponse.Vessels.PageInfo.HasNextPage, it.hasNext)
		is.Equal(mockResponse.Vessels.PageInfo.EndCursor, it.cursor)
		is.Equal([]byte(mockResponse.Vessels.PageInfo.EndCursor), it.position)
	})

	t.Run("loadBatch_Error", func(t *testing.T) {
		is := is.New(t)
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)
		is.NoErr(err)

		client.RunFn = func(ctx context.Context, req *graphql.Request, resp interface{}) error {
			return errors.New("some error")
		}

		err = it.loadBatch(context.Background())

		is.True(err != nil)
		is.Equal(err.Error(), "error making graphQL Request: some error")
	})

	t.Run("loadBatch_Retry", func(t *testing.T) {
		is := is.New(t)
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		batchSize := 100
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, batchSize, position)
		is.NoErr(err)

		// Mock the GraphQL response
		mockResponse := struct {
			Vessels Vessels
		}{
			Vessels: Vessels{
				Nodes: []Node{},
				PageInfo: PageInfo{
					HasNextPage: false,
					EndCursor:   "some_cursor",
				},
			},
		}

		retries := 0
		maxRetries := 2

		client.RunFn = func(ctx context.Context, req *graphql.Request, resp interface{}) error {
			if retries < maxRetries {
				retries++
				return errors.New("some error")
			}
			arg := resp.(*struct{ Vessels Vessels })
			*arg = mockResponse
			return nil
		}

		err = it.loadBatch(context.Background())

		is.NoErr(err)
		is.Equal(mockResponse.Vessels.Nodes, it.currentBatch)
		is.Equal(mockResponse.Vessels.PageInfo.HasNextPage, it.hasNext)
		is.Equal(mockResponse.Vessels.PageInfo.EndCursor, it.cursor)
		is.Equal([]byte(mockResponse.Vessels.PageInfo.EndCursor), it.position)
		is.Equal(retries, maxRetries) // Check if the retries have been exhausted
	})
}
