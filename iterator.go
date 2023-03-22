package ais

import (
	"context"
	"encoding/json"
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
	sdk.Logger(ctx).Info().Msg("in Next()")
	// return next message from cached batch
	var out Node
	if len(it.currentBatch) > 0 {
		sdk.Logger(ctx).Info().Msg("currentBatch not empty")
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	} else {
		err := it.loadBatch(ctx)
		if err != nil {
			sdk.Logger(ctx).Info().Msgf("loadBatch returned error %w", err)
			return sdk.Record{}, err
		}
		out, it.currentBatch = it.currentBatch[0], it.currentBatch[1:]
	}

	sdk.Logger(ctx).Info().Msg("returning record")
	return wrapAsRecord(out, it.position)
}

func (it *Iterator) loadBatch(ctx context.Context) error {
	sdk.Logger(ctx).Info().Msg("loadBatch called")
	graphqlRequest := graphql.NewRequest(it.query)
	graphqlRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", it.token))
	graphqlRequest.Var("first", 100)
	var Response struct {
		Vessels Vessels
	}
	if it.hasNext {
		sdk.Logger(ctx).Info().Msg("hasNext is true")
		graphqlRequest.Var("after", it.cursor)
	}
	sdk.Logger(ctx).Info().Msg("calling graphQL endpoint")
	if err := it.client.Run(context.Background(), graphqlRequest, &Response); err != nil {
		sdk.Logger(ctx).Err(err).Msgf("graphqlRequest: %+v", graphqlRequest)
		return err
	}
	it.currentBatch = Response.Vessels.Nodes
	it.hasNext = Response.Vessels.PageInfo.HasNextPage
	it.cursor = Response.Vessels.PageInfo.EndCursor
	it.position = []byte(Response.Vessels.PageInfo.EndCursor)

	return nil
}

func wrapAsRecord(in Node, endCursor sdk.Position) (sdk.Record, error) {
	sdk.Logger(context.Background()).Info().Msgf("record: %+v", in)
	updateTimestamp, err := time.Parse(time.RFC3339, in.UpdateTimestamp)
	if err != nil {
		sdk.Logger(context.Background()).Info().Msgf("timestamp: %s; error %w", in.UpdateTimestamp, err)
		return sdk.Record{}, err
	}

	sdkMetadata := make(sdk.Metadata)
	sdkMetadata.SetCreatedAt(updateTimestamp)

	b, err := json.Marshal(in)
	if err != nil {
		return sdk.Record{}, err
	}

	return sdk.Util.Source.NewRecordCreate(endCursor, sdkMetadata, nil, sdk.RawData(b)), nil
}
