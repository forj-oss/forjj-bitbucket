// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/forj-oss/goforjj"
)

// DoCreate --> plugin tasks for forjj create
// req_data contains the request data posted by forjj. Structure generated from 'bitbucket.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {

	//Get instance name
	instance := req.Forj.ForjjInstanceName

	//init BitbucketPlugin
	bbs := BitbucketPlugin{
		sourcePath:  req.Forj.ForjjSourceMount,
		deployMount: req.Forj.ForjjDeployMount,
		instance:    req.Forj.ForjjInstanceName,
		deployTo:    req.Forj.ForjjDeploymentEnv,
		key:         req.Objects.App[instance].Key,
		secret:      req.Objects.App[instance].Secret,
		team:        req.Objects.App[instance].Team,
	}

	log.Printf("Checking parameters : %#v", bbs)

	//Check key, secret and source
	check := make(map[string]bool)
	check["key"] = true
	check["secret"] = true
	check["source"] = true

	//ensure source path is writeable && key/secret isn't empty
	if bbs.verifyReqFails(ret, check) {
		return
	}

	//verify if instance existe
	if a, found := req.Objects.App[instance]; !found {
		ret.Errorf("Internal issue. Forjj has not given the Application information for '%s'. Aborted.")
		return
	} else {
		bbs.app = &a
	}

	//Check bitbucket connection
	log.Println("Checking bitbucket connection.")
	ret.StatusAdd("Connect to bitbucket...")

	if Bclient := bbs.bitbucketConnect("X", ret); Bclient == nil { //!\\
		return
	}

	// Init Team
	if !req.InitTeam(&bbs) {
		ret.Errorf("Internal Error. Unable to define the team.")
	}

	//Init Project
	if !req.InitProject(&bbs) {
		ret.Errorf("Internal Error. Unable to define project.")
	}

	//Create yaml data for maintain function
	if err := bbs.createYamlData(req, ret); err != nil {
		ret.Errorf("Unable to create. %s", err)
		return
	}

	//repos exist ?
	if err := bbs.reposExists(ret); err != nil {
		ret.Errorf("%s\nUnable to 'create' your forge when gitlab already has an infra project created. Clone it and use 'update' instead.", err)
		return 419
	}

	//CheckSourceExistence
	if err := bbs.checkSourcesExistence("create"); err != nil {
		ret.Errorf("%s\nUnable to 'create' your forge", err)
		return
	}
	ret.StatusAdd("Environment checked. Ready to be created.")

	//Path in ctxt git
	gitFile := path.Join(bbs.instance, bitbucketFile)

	//Save bitbucket source
	if _, err := bbs.saveYaml(&bbs.bitbucketSource, bbs.sourceFile); err != nil {
		ret.Errorf("%s", err)
		return
	}
	log.Printf(ret.StatusAdd("Configuration saved in source repo '%s' (%s).", gitFile, bbs.sourcePath))

	//Save bitbucket deploy
	if _, err := bbs.saveYaml(&bbs.bitbucketDeploy, bbs.deployFile); err != nil {
		ret.Errorf("%s", err)
		return
	}
	log.Printf(ret.StatusAdd("Configuration saved in deploy repo '%s' (%s).", gitFile, path.Join(bbs.deployMount, bbs.deployTo)))

	//Build final post answer
	for k, v := range bbs.bitbucketSource.Urls {
		ret.Services.Urls[k] = v
	}

	//API by forjj
	ret.Services.Urls["api_url"] = bbs.bitbucketSource.Urls["bitbucket-base-url"]

	ret.CommitMessage = fmt.Sprint("Bitbucket configuration created.")

	ret.AddFile(goforjj.FilesSource, gitFile)
	ret.AddFile(goforjj.FilesDeploy, gitFile)

	ret.StatusAdd("end")
	return
}

// DoUpdate --> plugin task for forjj update
// req_data contains the request data posted by forjj. Structure generated from 'bitbucket.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
	//TODO
	return
}

// DoMaintain --> plugin task for forjj maintain
// req_data contains the request data posted by forjj. Structure generated from 'bitbucket.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {
	instance := req.Forj.ForjjInstanceName

	var bbs BitbucketPlugin
	if a, found := req.Objects.App[instance]; !found {
		ret.Errorf("Invalid request. Missing Objects/App/%s", instance)
		return
	} else {
		bbs = BitbucketPlugin{
			deployMount:    req.Forj.ForjjDeployMount,
			workspaceMount: req.Forj.ForjjWorkspaceMount,
			key:            a.Key,
			secret:         a.Secret,
			maintainCtxt:   true,
			force:          req.Forj.Force == "true",
			sourcePath:     req.Forj.ForjjSourceMount,
		}
	}

	check := make(map[string]bool)
	check["token"] = true
	check["workspace"] = true
	check["deploy"] = true

	if bbs.verifyReqFails(ret, check) {
		return
	}

	confFile := path.Join(bbs.deployMount, req.Forj.ForjjDeploymentEnv, instance, bitbucketFile)

	//read yaml file
	if err := bbs.loadYaml(confFile /*ret, instance*/); err != nil {
		ret.Errorf("%s" /*, err*/)
		return
	}

	if bbs.bitbucketConnect("", ret) == nil {
		return
	}

	if !bbs.ensureTeamExists(ret) {
		return
	}

	if !bbs.IsNewForge(ret) {
		return
	}

	if bbs.bitbucketDeploy.NoRepos {
		log.Printf(ret.StatusAdd("Repos maintained limited to your infra project"))
	}

	//loop verif
	for name, repoData := range bbs.bitbucketDeploy.Repos {
		if !repoData.Infra && bbs.bitbucketDeploy.NoRepos {
			log.Printf(ret.StatusAdd("Project ignored: %s", name))
			continue
		}
		if repoData.Role == "infra" && !repoData.IsDeployable {
			log.Printf(ret.StatusAdd("Project ignored: %s - Infra project owned by '%s'", name, bbs.bitbucketDeploy.ProdTeam))
			continue
		}
		if err := repoData.ensureExists(&bbs, ret); err != nil {
			return
		}
		//...
		log.Printf(ret.StatusAdd("Project maintained: %s", name))
	}

	return
}
