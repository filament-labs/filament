package domain

import (
	"github.com/codemaestro64/filament/apps/api/pkg/util"
)

type Settings struct {
	Network util.Network `json:"network"`
}

type SettingsRequest struct{}
type SettingsResponse struct {
}
