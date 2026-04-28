package settings

import (
	"github.com/teadove/teasutils/service_utils/settings_utils"
)

type baseSettings struct {
	AIURL    string `env:"AI_URL"`
	AIAPIKEY string `env:"AI_API_KEY"`
}

// Settings
// nolint: gochecknoglobals // need it
var Settings = settings_utils.MustGetSetting[baseSettings]("VIBESITER_")
