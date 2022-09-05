package handler

import (
	"fmt"
	"testing"
)

func TestMeta(t *testing.T) {
	user, pkg := "kaiya", "goutils"
	importPath := fmt.Sprintf("%s/%s/%s", importHost, user, pkg)
	meta := fmt.Sprintf(formatStr, Hostname, importPath)
	t.Logf("%s", meta)
}
