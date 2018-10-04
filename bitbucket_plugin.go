// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"fmt"
)

type BitbucketPlugin struct {
	yaml        YamlBitbucket
	source_path string

	//sourceMount 	string
	deployMount 	string
	instance 		string
	deployTo 		string
	key 			string
	secret 			string
	team			string

	app 			*AppInstanceStruct  	//permet l'accès au forjfile
	Client 			*bitbucket.Client 			//client bitbucket cf api bitbucket 
	bitbucket_source 	BitbucketSourceStruct 		//urls...
	bitbucketDeploy 	BitbucketDeployStruct     	//

	gitFile 		string
	deployFile 		string
	sourceFile  	string
}

const bitbucket_file = "forjj-bitbucket.yaml"

type YamlBitbucket struct {
}

func new_plugin(src string) (p *BitbucketPlugin) {
	p = new(BitbucketPlugin)

	p.source_path = src
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

func (p *BitbucketPlugin) save_yaml(in interface{}, file string) (Updated bool, _ error){
	d, err := yaml.Marshal(in)
	if err != nil {
		return false, fmt.Errorf("Unable to encode bitbucket data in yaml. %s", err)
	}

	if d_before, err := ioutil.ReadFile(file); err != nil{
		Updated = true
	} else {
		Updated = (string(d) != string(d_before))
	}
	
	if !Updated {
		return
	}
	if err = ioutil.WriteFile(file, d, 0644); err != nil {
		return false, fmt.Errorf("Unable to save '%s'. %s", file, err)
	}
	return
}

func (p *BitbucketPlugin) load_yaml(ret *goforjj.PluginData, instance string) (status bool) {
	file := path.Join(instance, bitbucket_file)

	d, err := ioutil.ReadFile(file)
	if err != nil {
		ret.Errorf("Unable to load '%s'. %s", file, err)
		return
	}

	err = yaml.Unmarshal(d, &p.yaml)
	if err != nil {
		ret.Errorf("Unable to decode forjj bitbucket data in yaml. %s", err)
		return
	}
	return true
}
