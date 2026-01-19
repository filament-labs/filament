package filwallet

import (
	"errors"
	"fmt"
	"time"
)

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
)

func (e Env) IsProduction() bool {
	return e == Production
}

func (e Env) Validate() error {
	switch e {
	case Development, Production:
		return nil
	default:
		return fmt.Errorf("invalid environment: %s", e)
	}
}

type Config struct {
	Env             Env
	SessionDuration time.Duration
	RPCEndpoint     string
	RPCToken        string
	DataDir         string
}

func (c *Config) Validate() error {
	if err := c.Env.Validate(); err != nil {
		return err
	}

	if c.DataDir == "" {
		return errors.New("missing app data directory")
	}

	if c.RPCEndpoint == "" {
		return errors.New("missing api endpoint")
	}

	return nil
}
