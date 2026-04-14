package utils

import (
	"errors"

	"github.com/odysseythink/confy"
	"mlib.com/gofy/server/main/sandbox/runner/types"
)

var (
	ErrNetworkDisabled = errors.New("network is disabled, please enable it in the configuration")
)

func CheckOptions(options *types.RunnerOptions) error {
	if options.EnableNetwork && !confy.GetWithDefault[bool]("enable_network", false) {
		return ErrNetworkDisabled
	}

	return nil
}
