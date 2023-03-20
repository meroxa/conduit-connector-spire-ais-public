package ais

import (
	"fmt"
)

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// GlobalConfigParam is named global_config_param_name and needs to be
	// provided by the user.
	apiUrl string
	token  string
	query  string
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}
