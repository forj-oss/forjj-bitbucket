// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
)

//CheckSourceExistence ...
func (r *UpdateReq) CheckSourceExistence(ret *goforjj.PluginData) (p *BitbucketPlugin, status bool) {
	log.Print("Checking Bitbucket source code existence.")
	srcPath := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(srcPath, bitbucketFile)); err == nil {
		log.Printf(ret.Errorf("Unable to create the bitbucket source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update bitbucket according to his configuration. %s.", srcPath, srcPath, err))
		return
	}

	p = new_plugin(srcPath)

	ret.StatusAdd("environment checked.")
	return p, true
}

//SaveMaintainOptions function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *UpdateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) goforjj.PluginOption {
	opt, ok := options[option]
	if ok && value != "" {
		opt.Value = value
		return opt
	}
	if !ok {
		opt = goforjj.PluginOption{Help: help}
		if value == "" {
			opt.Value = defaultv
		} else {
			opt.Value = value
		}
	}
	return opt
}

//SetRepo TODO
func (bbs *BitbucketPlugin) SetRepo(repo *RepoInstanceStruct, isInfra, isDeployable bool) {
	upstream := bbs.DefineRepoUrls(repo.Name)

	owner := bbs.bitbucketDeploy.Team
	if isInfra {
		owner = bbs.bitbucketDeploy.ProdTeam
	}

	//set it, found or not
	pjt := RepositoryStruct{}
	pjt.set(repo,
				map[string]goforjj.PluginRepoRemoteUrl{"origin": upstream},
				map[string]string{"master": "origin/master"},
				isInfra,
				isDeployable, owner)
	bbs.bitbucketDeploy.Repos[repo.Name] = pjt
}
