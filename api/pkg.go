package handler

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	formatStr  = `<!DOCTYPE html><html lang="en"><head><meta name="go-import" content="%s git %s" /></head></html>`
	Hostname   = "kaiya.js.org"
	importHost = "https://github.com"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	isGoGet := r.FormValue("go-get")
	if isGoGet == "1" {
		// go-get request
		fmt.Fprintf(w, "from go get")
	} else {

	}
	user, pkg, version, onlyPkg := splitPkgName(r.URL.Path)
	remotePath := fmt.Sprintf("%s/%s/%s", importHost, user, pkg)
	var importPrefix string
	if onlyPkg {
		importPrefix = fmt.Sprintf("%s/%s", Hostname, pkg)
	} else {
		importPrefix = fmt.Sprintf("%s/%s/%s", Hostname, user, pkg)
	}

	fmt.Fprintf(w, formatStr, importPrefix, remotePath)
	w.Write([]byte(fmt.Sprintf("user:%s, pkg:%s, version:%s", user, pkg, version)))
}

func splitPkgName(originalPath string) (user, pkg, version string, onlyPkg bool) {
	// support two case:
	// 1. example.org/goutils.v1
	// 2. exmaple.org/kaiya/goutils.v1
	path := strings.TrimPrefix(originalPath, "/")
	var pkgWithVer string
	tmpPath := strings.SplitN(path, "/", 2)
	if len(tmpPath) == 1 {
		onlyPkg = true
		pkgWithVer = tmpPath[0]
	} else if len(tmpPath) == 2 {
		user, pkgWithVer = tmpPath[0], tmpPath[1]
	}
	// split version
	tmpPkg := strings.SplitN(pkgWithVer, ".", 2)
	pkg = tmpPkg[0]
	version = "v1"
	if len(tmpPkg) == 2 {
		version = tmpPkg[1]
	}
	if onlyPkg {
		user = fmt.Sprintf("go-%s", pkg)
	}
	return
}
