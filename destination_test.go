package ais_test

import (
	"context"
	"testing"

	ais "github.com/meroxa/conduit-connector-spire-ais"
)

func TestTeardown_NoOpen(t *testing.T) {
	con := ais.NewDestination()
	err := con.Teardown(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
