// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type BitbucketPlugin struct {
	yaml        YamlBitbucket
	sourcePath string

	//sourceMount 	string
	deployMount 	string
	instance 		string
	deployTo 		string
	key 			string
	secret 			string
	team			string

	app 			*AppInstanceStruct  	//permet l'accès au forjfile
	Client 			*bitbucket.Client 			//client bitbucket cf api bitbucket 
	bitbucketSource 	BitbucketSourceStruct 		//urls...
	bitbucketDeploy 	BitbucketDeployStruct     	//

	gitFile 		string
	deployFile 		string
	sourceFile  	string

	//maintain
	workspaceMount		string
	maintainCtxt		bool
	force			bool

	newForge		bool
}

type BitbucketDeployStruct struct{
	goforjj.PluginService 	`yaml:",inline"` //urls
	Repos 					map[string]RepositoryStruct // projects managed in gitlab
	NoRepos 				bool `yaml:",omitempty"`
	ProdTeam 				string
	Team 					string
	TeamDisplayName 		string
	//...
}

type BitbucketSourceStruct struct{
	goforjj.PluginService `,inline` //base url
	ProdTeam string `yaml:"production-team-name"`//`yaml:"production-group-name, omitempty"`
}

const bitbucketFile = "forjj-bitbucket.yaml"

type YamlBitbucket struct {
}

func new_plugin(src string) (p *BitbucketPlugin) {
	p = new(BitbucketPlugin)

	p.sourcePath = src
	return
}

func (p *BitbucketPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *BitbucketPlugin) load_from(ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *BitbucketPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *BitbucketPlugin) saveYaml(in interface{}, file string) (Updated bool, _ error){
	d, err := yaml.Marshal(in)
	if err != nil {
		return false, fmt.Errorf("Unable to encode bitbucket data in yaml. %s", err)
	}

	if dBefore, err := ioutil.ReadFile(file); err != nil{
		Updated = true
	} else {
		Updated = (string(d) != string(dBefore))
	}
	
	if !Updated {
		return
	}
	if err = ioutil.WriteFile(file, d, 0644); err != nil {
		return false, fmt.Errorf("Unable to save '%s'. %s", file, err)
	}
	return
}

func (p *BitbucketPlugin) loadYaml(file string) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Unable to load '%s'. %s", file, err)
	}

	err = yaml.Unmarshal(d, &p.bitbucketDeploy)

	if err != nil {
		return fmt.Errorf("Unable to decode bitbucket data in yaml. %s", err)
	}

	return nil
}
