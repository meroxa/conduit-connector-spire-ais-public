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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockGraphQLClient struct {
	mock.Mock
}

func (m *MockGraphQLClient) Run(ctx context.Context, req *graphql.Request, resp interface{}) error {
	args := m.Called(ctx, req, resp)
	return args.Error(0)
}

func TestIterator(t *testing.T) {
	t.Run("NewIterator", func(t *testing.T) {
		client := &graphql.Client{}
		token := "test-token"
		query := "test-query"
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, position)

		require.NoError(t, err)
		assert.Equal(t, client, it.client)
		assert.Equal(t, token, it.token)
		assert.Equal(t, query, it.query)
		assert.Equal(t, position, it.position)
	})

	t.Run("HasNext", func(t *testing.T) {
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, position)
		require.NoError(t, err)

		// Set up expected behavior
		expectedResponse := struct {
			Vessels Vessels
		}{
			Vessels: Vessels{
				PageInfo: PageInfo{HasNextPage: false, EndCursor: ""},
				Nodes:    []Node{},
			},
		}
		client.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(2).(*struct {
				Vessels Vessels
			})
			*arg = expectedResponse
		}).Once()

		// hasNext is true
		it.hasNext = true
		assert.True(t, it.HasNext(context.Background()))

		// hasNext is false, but there are more nodes
		it.hasNext = false
		it.currentBatch = []Node{{}, {}}
		assert.True(t, it.HasNext(context.Background()))

		// hasNext is false, and no more nodes
		it.currentBatch = []Node{}
		assert.False(t, it.HasNext(context.Background()))
	})

	t.Run("Next", func(t *testing.T) {
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, position)
		require.NoError(t, err)

		it.currentBatch = []Node{
			{
				UpdateTimestamp: "2021-10-01T15:00:00Z",
			},
		}

		record, err := it.Next(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, record)
	})

	t.Run("loadBatch", func(t *testing.T) {
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, position)
		require.NoError(t, err)

		client.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		err = it.loadBatch(context.Background())
		require.NoError(t, err)

		client.AssertExpectations(t)
	})

	t.Run("loadBatchError", func(t *testing.T) {
		client := &MockGraphQLClient{}
		token := "test-token"
		query := "test-query"
		position := sdk.Position("test-position")

		it, err := NewIterator(client, token, query, position)
		require.NoError(t, err)

		client.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphql error")).Once()

		err = it.loadBatch(context.Background())
		require.Error(t, err)

		client.AssertExpectations(t)
	})
}
