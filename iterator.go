package ais

import (
	"context"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/machinebox/graphql"
)

type Iterator struct {
	token        string
	query        string
	cursor       string
	currentBatch []Node
	client       *graphql.Client
	position     sdk.Position
}

func NewIterator(client *graphql.Client, token string, query string, p sdk.Position) (*Iterator, error) {
	return &Iterator{
		token:    token,
		query:    query,
		client:   client,
		position: p,
	}, nil

}

func (it *Iterator) HasNext(ctx context.Context) (bool, error) {

	// return early if there are more nodes
	if len(it.currentBatch) > 0 {
		return true, nil
	}

	// there are additional pages
	if it.cursor != "" {
		return true, nil
	}

	// we don't have any more nodes and we don't know if there are additional pages
	// loadBatch returns sdk.ErrBackoffRetry if there are no more pages
	if err := it.loadBatch(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (it *Iterator) loadBatch(ctx context.Context) error {
	graphqlRequest := graphql.NewRequest(it.query)
	graphqlRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))
}
