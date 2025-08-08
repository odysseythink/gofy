package utils

import (
	"errors"

	"github.com/spf13/viper"
	"mlib.com/gofy/server/main/sandbox/runner/types"
)

var (
	ErrNetworkDisabled = errors.New("network is disabled, please enable it in the configuration")
)

func CheckOptions(options *types.RunnerOptions) error {
	if options.EnableNetwork && !viper.GetBoolWithDefault("enable_network", false) {
		return ErrNetworkDisabled
	}

	return nil
}
