package filwallet

import (
	"errors"

	"github.com/codemaestro64/filament/apps/api/pkg/util"
)

type Config struct {
	Network        util.Network
	SessionTimeout int64
	RPCEndpoint    string
	RPCToken       string
	DataDir        string
}

func (c *Config) Validate() error {
	if c.DataDir == "" {
		return errors.New("missing app data directory")
	}

	if c.RPCEndpoint == "" {
		return errors.New("missing api endpoint")
	}

	return nil
}
