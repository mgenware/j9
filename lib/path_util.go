package lib

import (
	"os"
	"strings"

	"github.com/mgenware/goutil/stringsx"
)

const homeEnv = "$HOME"

func FormatPath(s string, evaluate bool) string {
	if strings.HasPrefix(s, "~/") {
		s = homeEnv + "/" + stringsx.SubStringFromStart(s, 2)
	} else if s == "~" {
		s = homeEnv
	}
	if evaluate {
		return os.ExpandEnv(s)
	}
	return s
}
