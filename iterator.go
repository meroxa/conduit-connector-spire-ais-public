package ais

import (
	"context"
	"fmt"
	"time"

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

func (it *Iterator) HasNext(ctx context.Context) bool {

	// return early if there are more nodes
	if len(it.currentBatch) > 0 {
		return true
	}

	// there are additional pages
	if it.cursor != "" {
		return true
	}

	// we don't have any more nodes and we don't know if there are additional pages
	// loadBatch returns sdk.ErrBackoffRetry if there are no more pages
	if err := it.loadBatch(ctx); err != nil {
		return false
	}

	return true
}

func (it *Iterator) Next(ctx context.Context) (sdk.Record, error) {
	// return next message from cached batch
	var out Node
	if len(it.currentBatch) > 0 {
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	}

	return wrapAsRecord(out)
}

func (it *Iterator) loadBatch(ctx context.Context) error {
	graphqlRequest := graphql.NewRequest(it.query)
	graphqlRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))

	return nil
}

func wrapAsRecord(in Node) (sdk.Record, error) {
	// todo: use endCursor for position
	//position := in.UpdateTimestamp

	updateTimestamp, err := time.Parse(time.RFC3339, in.UpdateTimestamp)
	if err != nil {
		return sdk.Record{}, err
	}

	sdkMetadata := make(sdk.Metadata)
	sdkMetadata.SetCreatedAt(updateTimestamp)

	return sdk.Util.Source.NewRecordCreate(nil, sdkMetadata, nil, in.toStructuredData()), nil
}
