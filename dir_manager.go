package j9

import (
	"os"
	"strings"

	"github.com/mgenware/goutil/stringsx"
	"github.com/mgenware/j9/lib"
)

type dirManager struct {
	lastDir string
}

func (d *dirManager) LastDir() string {
	return d.lastDir
}

func (d *dirManager) NextWD(cmd string, expandLocally bool) string {
	if strings.HasPrefix(cmd, "cd") {
		var dir string
		if len(cmd) == 2 {
			dir = os.Getenv("HOME")
		} else if cmd[2] == ' ' && len(cmd) > 3 {
			dir = strings.TrimSpace(stringsx.SubStringFromStart(cmd, 3))
		}

		if dir != "" {
			dir = lib.FormatPath(dir, expandLocally)
			d.lastDir = dir
			return dir
		}
	}
	return ""
}
