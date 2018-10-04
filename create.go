// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
	"fmt"
)

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *BitbucketPlugin, status bool) {
	log.Print("Checking Bitbucket source code existence.")
	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src_path, bitbucket_file)); err == nil {
		log.Printf(ret.Errorf("Unable to create the bitbucket source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update bitbucket according to his configuration. %s.", src_path, src_path, err))
		return
	}

	p = new_plugin(src_path)

	log.Printf(ret.StatusAdd("environment checked."))
	return p, true
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func (bbs *BitbucketPlugin) create_yaml_data(req *CreateReq, ret *goforjj.PluginData) error{
	if bbs.bitbucket_source.Urls == nil{
		return fmt.Errorf("Internal Error. Urls was not set")
	}

	bbs.bitbucketDeploy.Repos = make(map[string]RepositoryStruct)
	//bbs.bitbucketDeploy.Users = ...

	//Norepo
	bbs.bitbucketDeploy.NoRepos = (bbs.app.ReposDisabled == "true")
	if bbs.bitbucketDeploy.NoRepos {
		log.Print("Repositories_disabled is true. forjj_gitlab won't manage repositories except the infra repository.")
	}

	//SetOrgHooks

	for name, repo := range req.Objects.Repo{
		is_infra := (name == bbs.app.ForjjInfra)
		if bbs.bitbucketDeploy.NoRepos && !is_infra {
			continue
		}
		if !repo.IsValid(name, ret){
			ret.StatusAdd("Warning!!! Invalid repo '%s' requested. Ignored.")
			continue
		}
		bbs.SetRepo(&repo, is_infra, repo.Deployable == "true")
		//bbs.SetHooks(...)

	}

	log.Printf("forjj-bitbucket manages %d repo(s).", len(bbs.bitbucketDeploy.Repos))

	//more ...

	return nil
}

// DefineRepoUrls return default repo url for the repo name given
func (bbs *BitbucketPlugin) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl){
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: bbs.bitbucket_source.Urls["bitbucket-ssh"] + bbs.bitbucketDeploy.Team + "/" + name + ".git",
		Url: bbs.bitbucket_source.Urls["bitbucket-url"] + "/" + bbs.bitbucketDeploy.Team + "/" + name,
	}
	return
}
