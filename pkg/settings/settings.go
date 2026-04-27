package settings

import (
	"github.com/teadove/teasutils/service_utils/settings_utils"
)

type baseSettings struct{}

// Settings
// nolint: gochecknoglobals // need it
var Settings = settings_utils.MustGetSetting[baseSettings]("VIBESITER_")
