package main

import (
	"fmt"
	"os"

	"github.com/forj-oss/goforjj"
	"golang.org/x/sys/unix"
)

//verifyReqFails ...
func (bbs *BitbucketPlugin) verifyReqFails(ret *goforjj.PluginData, check map[string]bool) bool {
	if v, ok := check["source"]; ok && v {
		if reqCheckPath("source (forjj-source-mount)", bbs.sourcePath, ret) {
			return true
		}
	}

	if v, ok := check["key"]; ok && v {
		if bbs.key == "" {
			ret.ErrorMessage = fmt.Sprint("bitbucket key is empty - Required")
			return true
		}
	}

	if v, ok := check["secret"]; ok && v {
		if bbs.secret == "" {
			ret.ErrorMessage = fmt.Sprint("bitbucket secret is empty - Required")
			return true
		}
	}

	return false
}

//reqCheckPath check path is writable.
// return false if something is wrong.
func reqCheckPath(name, path string, ret *goforjj.PluginData) bool {

	if path == "" {
		ret.ErrorMessage = name + " is empty."
		return true
	}

	if _, err := os.Stat(path); err != nil {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is inexistent.", path)
		return true
	}

	if !IsWritable(path) {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is NOT writable", path)
		return true
	}

	return false
}

//IsWritable Linux support only
func IsWritable(path string) (res bool) {
	return unix.Access(path, unix.W_OK) == nil
}

//
func inStringList(element string, elements ...string) string {
	for _, value := range elements {
		if element == value {
			return value
		}
	}
	return ""
}
