// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"log"
	"os"

	"github.com/forj-oss/goforjj"
)

//CheckSourceExistence Return ok if the jenkins instance exist
func (r *MaintainReq) CheckSourceExistence(ret *goforjj.PluginData) (status bool) {
	log.Print("Checking bitbucket source code path existence.")

	if _, err := os.Stat(r.Forj.ForjjSourceMount); err == nil {
		ret.Errorf("Unable to maintain bitbucket instances. '%s' is inexistent or innacessible.\n",
			r.Forj.ForjjSourceMount)
		return
	}
	ret.StatusAdd("environment checked.")
	return true
}

//instantiate ...
func (r *MaintainReq) instantiate(ret *goforjj.PluginData) (status bool) {

	return true
}

//MaintainTeamHooks TODO
func (bbs *BitbucketPlugin) MaintainTeamHooks(ret *goforjj.PluginData) (_ bool) {
	return
}
