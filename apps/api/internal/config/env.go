package config

import "fmt"

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
	Staging     Env = "staging"
	Test        Env = "test"
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
