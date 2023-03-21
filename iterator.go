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
	hasNext      bool
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

	if it.hasNext {
		it.loadBatch(ctx)
		return true
	}

	return false
}

func (it *Iterator) Next(ctx context.Context) (sdk.Record, error) {
	// return next message from cached batch
	var out Node
	if len(it.currentBatch) > 0 {
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	}

	return wrapAsRecord(out, it.position)
}

func (it *Iterator) loadBatch(ctx context.Context) error {
	graphqlRequest := graphql.NewRequest(it.query)
	graphqlRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))
	var response Vessels
	if it.hasNext {
		graphqlRequest.Var("after", it.cursor)
	}
	if err := it.client.Run(context.Background(), graphqlRequest, &response); err != nil {
		panic(err)
	}
	it.currentBatch = response.Nodes
	it.hasNext = response.PageInfo.HasNextPage
	it.cursor = response.PageInfo.EndCursor
	it.position = []byte(response.PageInfo.EndCursor)

	return nil
}

func wrapAsRecord(in Node, endCursor sdk.Position) (sdk.Record, error) {
	updateTimestamp, err := time.Parse(time.RFC3339, in.UpdateTimestamp)
	if err != nil {
		return sdk.Record{}, err
	}

	sdkMetadata := make(sdk.Metadata)
	sdkMetadata.SetCreatedAt(updateTimestamp)

	return sdk.Util.Source.NewRecordCreate(endCursor, sdkMetadata, nil, in.toStructuredData()), nil
}
