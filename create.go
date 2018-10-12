// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/forj-oss/goforjj"
)

//CheckSourceExistence return true if instance doesn't exist.
func (r *CreateReq) CheckSourceExistence(ret *goforjj.PluginData) (p *BitbucketPlugin, status bool) {
	log.Print("Checking Bitbucket source code existence.")
	srcPath := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(srcPath, bitbucketFile)); err == nil {
		log.Printf(ret.Errorf("Unable to create the bitbucket source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update bitbucket according to his configuration. %s.", srcPath, srcPath, err))
		return
	}

	p = new_plugin(srcPath)

	log.Printf(ret.StatusAdd("environment checked."))
	return p, true
}

//SaveMaintainOptions ...
func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

//createYamlData ...
func (bbs *BitbucketPlugin) createYamlData(req *CreateReq, ret *goforjj.PluginData) error {
	if bbs.bitbucketSource.Urls == nil {
		return fmt.Errorf("Internal Error. Urls was not set")
	}

	bbs.bitbucketDeploy.Repos = make(map[string]RepositoryStruct)
	//bbs.bitbucketDeploy.Users = ...

	//Norepo
	bbs.bitbucketDeploy.NoRepos = (bbs.app.ReposDisabled == "true")
	if bbs.bitbucketDeploy.NoRepos {
		log.Print("Repositories_disabled is true. forjj_gitlab won't manage repositories except the infra repository.")
	}

	//SetTeamHooks
	bbs.setTeamHooks(bbs.app.TeamWebhooksDisabled, bbs.app.ReposWebhooksDisabled, bbs.app.TeamHookPolicy, req.Objects.Webhooks)

	for name, repo := range req.Objects.Repo {
		isInfra := (name == bbs.app.ForjjInfra)
		if bbs.bitbucketDeploy.NoRepos && !isInfra {
			continue
		}
		if !repo.IsValid(name, ret) {
			ret.StatusAdd("Warning!!! Invalid repo '%s' requested. Ignored.")
			continue
		}
		bbs.SetRepo(&repo, isInfra, repo.Deployable == "true")
		bbs.SetHooks(&repo, req.Objects.Webhooks)

	}

	log.Printf("forjj-bitbucket manages %d repo(s).", len(bbs.bitbucketDeploy.Repos))

	//more ...

	return nil
}

// DefineRepoUrls return default repo url for the repo name given
func (bbs *BitbucketPlugin) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl) {
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: bbs.bitbucketSource.Urls["bitbucket-ssh"] + bbs.bitbucketDeploy.Team + "/" + name + ".git",
		Url: bbs.bitbucketSource.Urls["bitbucket-url"] + "/" + bbs.bitbucketDeploy.Team + "/" + name,
	}
	return
}
