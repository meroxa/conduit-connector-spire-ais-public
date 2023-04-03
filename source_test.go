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
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
	"github.com/stretchr/testify/mock"
)

type MockIteratorCreator struct {
	mock.Mock
	HasNext bool
	Next    sdk.Record
}

func (m *MockIteratorCreator) NewIterator(client GraphQLClient, token string, query string, batchSize int, p sdk.Position) (*Iterator, error) {
	args := m.Called(client, token, query, p)
	return args.Get(0).(*Iterator), args.Error(1)
}

func TestSource(t *testing.T) {
	t.Run("Parameters", func(t *testing.T) {
		source := NewSource()
		params := source.Parameters()

		is := is.New(t)
		is.True(params != nil)
		is.Equal(params["apiUrl"].Default, "https://api.spire.com/graphql")
		is.Equal(params["query"].Default, "")
		is.Equal(params["token"].Default, "")
		is.Equal(params["batchSize"].Default, "100")
	})

	t.Run("Configure", func(t *testing.T) {
		source := NewSource()
		cfg := map[string]string{
			"apiUrl":    "https://api.example.com/graphql",
			"token":     "test-token",
			"query":     "test-query",
			"batchSize": "100",
		}

		err := source.Configure(context.Background(), cfg)
		is := is.New(t)
		is.NoErr(err)
		is.Equal(source.(*Source).config.APIURL, "https://api.example.com/graphql")
		is.Equal(source.(*Source).config.Query, "test-query")
		is.Equal(source.(*Source).config.Token, "test-token")
		is.Equal(source.(*Source).config.BatchSize, 100)
	})

	t.Run("Open", func(t *testing.T) {
		source := NewSource()
		cfg := map[string]string{
			"apiUrl":     "https://api.example.com/graphql",
			"token":      "test-token",
			"query":      "test-query",
			"batch_size": "100",
		}

		err := source.Configure(context.Background(), cfg)
		is := is.New(t)
		is.NoErr(err)

		// Mock the iterator to be used in the Open method
		mockIterator := &Iterator{}
		mockIteratorCreator := &MockIteratorCreator{}
		mockIteratorCreator.On("NewIterator", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockIterator, nil).Once()

		source.(*Source).iteratorCreator = mockIteratorCreator

		err = source.Open(context.Background(), nil)
		is.NoErr(err)
		is.Equal(mockIterator, source.(*Source).iterator)

		mockIteratorCreator.AssertExpectations(t)
	})

	t.Run("Ack", func(t *testing.T) {
		source := NewSource()

		err := source.Ack(context.Background(), nil)
		is := is.New(t)
		is.NoErr(err)
	})

	t.Run("Teardown", func(t *testing.T) {
		source := NewSource()

		err := source.Teardown(context.Background())
		is := is.New(t)
		is.NoErr(err)
	})
}
