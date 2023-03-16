package ais

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// version is set during the build process with ldflags (see Makefile).
// Default version matches default from runtime/debug.
var version = "(devel)"

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "spire-ais",
		Summary:     "A source connector for getting data from the spire-ais GraphQL API",
		Description: "This connector should connect to the Spire-AIS GraphQL API using a Bearer Token.  It should allow sending a query to the API to fetch new data",
		Version:     version,
		Author:      "Meroxa, Inc.",
	}
}
