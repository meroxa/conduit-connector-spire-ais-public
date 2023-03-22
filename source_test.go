package ais

import (
	"context"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockIteratorCreator struct {
	mock.Mock
	HasNext bool
	Next    (sdk.Record)
}

func (m *MockIteratorCreator) NewIterator(client GraphQLClient, token string, query string, p sdk.Position) (*Iterator, error) {
	args := m.Called(client, token, query, p)
	return args.Get(0).(*Iterator), args.Error(1)
}

func TestSource(t *testing.T) {
	t.Run("Parameters", func(t *testing.T) {
		source := NewSource()
		params := source.Parameters()

		assert.NotNil(t, params)
		assert.Contains(t, params, "query")
		assert.Contains(t, params, "batch_size")
	})

	t.Run("Configure", func(t *testing.T) {
		source := NewSource()
		cfg := map[string]string{
			"api_url":    "https://api.example.com/graphql",
			"token":      "test-token",
			"query":      "test-query",
			"batch_size": "100",
		}

		err := source.Configure(context.Background(), cfg)
		require.NoError(t, err)
	})

	t.Run("Open", func(t *testing.T) {
		source := NewSource()
		cfg := map[string]string{
			"api_url":    "https://api.example.com/graphql",
			"token":      "test-token",
			"query":      "test-query",
			"batch_size": "100",
		}

		err := source.Configure(context.Background(), cfg)
		require.NoError(t, err)

		// Mock the iterator to be used in the Open method
		mockIterator := &Iterator{}
		mockIteratorCreator := &MockIteratorCreator{}
		mockIteratorCreator.On("NewIterator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockIterator, nil).Once()

		source.(*Source).iteratorCreator = mockIteratorCreator

		err = source.Open(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, mockIterator, source.(*Source).iterator)

		mockIteratorCreator.AssertExpectations(t)
	})

	t.Run("Ack", func(t *testing.T) {
		source := NewSource()

		err := source.Ack(context.Background(), nil)
		require.NoError(t, err)
	})

	t.Run("Teardown", func(t *testing.T) {
		source := NewSource()

		err := source.Teardown(context.Background())
		require.NoError(t, err)
	})
}
