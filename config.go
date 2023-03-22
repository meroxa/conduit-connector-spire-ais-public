package ais

import (
	"fmt"
)

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// APIURL is the Spire API URL to use for accessing the Spire GraphQL API.
	APIURL string `json:"api_url" default:"https://api.spire.com/graphql"`

	// Token is the access token to use when accessing the Spire GraphQL API.
	Token string `json:"token" validate:"required"`
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}
