package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	ais "github.com/meroxa/conduit-connector-spire-ais-public"
)

func main() {
	sdk.Serve(ais.Connector)
}
