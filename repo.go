package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/forj-oss/goforjj"
)

type RepositoryStruct struct {
	Name         string
	Flow         string            `yaml:",omitempty"`
	Description  string            `yaml:",omitempty"`
	Disabled     bool              `yaml:",omitempty"`
	IssueTracker bool              `yaml:"issue_tracker,omitempty"`
	Users        map[string]string `yaml:",omitempty"`
	Owner        string            `yaml:",omitempty"`
	//Groups

	exist         bool `yaml:",omitempty"`
	remotes       map[string]goforjj.PluginRepoRemoteUrl
	branchConnect map[string]string

	//maintain
	Infra        bool   `yaml:",omitempty"`
	Role         string `yaml:",omitempty"`
	IsDeployable bool

	//Webhooks
	WebHooks map[string]WebHookStruct `yaml:",omitempty"`
}

//IsValid ...
func (r *RepoInstanceStruct) IsValid(repoName string, ret *goforjj.PluginData) (valid bool) {
	if r.Name == "" {
		ret.Errorf("Invalid repo '%s'. Name is empty.", repoName)
		return
	}
	if r.Name != repoName {
		ret.Errorf("Invalid repo '%s'. Name must be equal to '%s'. But the repo name is set to '%s'.", repoName, repoName, r.Name)
		return
	}
	valid = true
	return
}

//set TODO
func (r *RepositoryStruct) set(repo *RepoInstanceStruct, remotes map[string]goforjj.PluginRepoRemoteUrl, branchConnect map[string]string, isInfra, IsDeployable bool, owner string) *RepositoryStruct {
	if r == nil {
		r = new(RepositoryStruct)
	}
	r.Name = repo.Name
	r.Description = repo.Title

	//issueTracker

	r.Flow = repo.Flow
	r.Infra = isInfra

	//addUsers
	//Groups

	r.remotes = remotes
	r.branchConnect = branchConnect

	//webHooks

	r.Role = repo.Role
	r.Owner = owner
	r.IsDeployable = IsDeployable

	return r
}

//
func (bbs *BitbucketPlugin) SetHooks(reqRepo *RepoInstanceStruct, hooks map[string]WebhooksInstanceStruct) {
	repo := bbs.bitbucketDeploy.Repos[reqRepo.Name]
	repo.WebHooks = make(map[string]WebHookStruct)

	if bbs.bitbucketDeploy.NoRepoHook {
		return
	}

	for name, hook := range hooks {
		if hook.Team == "true" {
			continue
		}
		if inStringList(repo.Name, strings.Split(hook.Repos, ",")...) == "" {
			continue
		}
		data := WebHookStruct{
			Url:         hook.Url,
			Events:      strings.Split(hook.Events, ","),
			Enabled:     hook.Enabled,
			ContentType: hook.PayloadFormat,
		}
		if v, err := strconv.ParseBool(hook.SslCheck); err == nil {
			data.SSLCheck = v
			log.Printf("SSL Check '%s' => %t", name, v)
		} else {
			log.Printf("SSL Check has an invalid boolean string representation '%s'. Ignored. SSL Check is set to true.", name)
			data.SSLCheck = true
		}

		repo.WebHooks[name] = data
		bbs.bitbucketDeploy.Repos[reqRepo.Name] = repo
	}
}
