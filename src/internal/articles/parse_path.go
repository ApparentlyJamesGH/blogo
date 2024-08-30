package articles

import (
	"path/filepath"
	"strings"
)

// Given a path, return the slug and extension
func ParseFilePath(fPath string) (string, string) {
	filenameWithExt := filepath.Base(fPath)
	extension := filepath.Ext(filenameWithExt)
	filename := strings.TrimSuffix(filenameWithExt, extension)
	return filename, extension
}
