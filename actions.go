// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"path" //path. ...
	//"io/ioutil" // files function

	"github.com/forj-oss/goforjj"
	//"github.com/ktrysmt/go-bitbucket"
	"golang.org/x/sys/unix"
	//"gopkg.in/yaml.v2" //use yaml. ...
)

/*type BitbucketPlugin struct{


}*/

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

type RepositoryStruct struct {
	Name 			string
	Flow 			string 				`yaml:",omitempty"`
	Description 	string 				`yaml:",omitempty"`
	Disabled 		bool				`yaml:",omitempty"`
	IssueTracker 	bool 				`yaml:"issue_tracker,omitempty"`
	Users 			map[string]string 	`yaml:",omitempty"`
	Owner 			string 				`yaml:",omitempty"`
	//Groups

	exist 			bool 				`yaml:",omitempty"`
	remotes 		map[string]goforjj.PluginRepoRemoteUrl
	branchConnect 	map[string]string 	
	//...
}

func (bbs *BitbucketPlugin) verify_req_fails(ret *goforjj.PluginData, check map[string]bool) bool{
	if v, ok := check["source"]; ok && v {
		if reqCheckPath("source (forjj-source-mount)", bbs.source_path, ret){
			return true
		}
	}

	if v, ok := check["key"]; ok && v {
		if bbs.key == ""{
			ret.ErrorMessage = fmt.Sprint("bitbucket key is empty - Required")
			return true
		}
	}

	if v, ok := check["secret"]; ok && v {
		if bbs.secret == ""{
			ret.ErrorMessage = fmt.Sprint("bitbucket secret is empty - Required")
			return true
		}
	}

	return false
}

// check path is writable.
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

// Linux support only
func IsWritable(path string) (res bool) {
	return unix.Access(path, unix.W_OK) == nil
}

/*func (bbs *BitbucketPlugin) bitbucket_connect(server string, ret *goforjj.PluginData) *bitbucket.Client{
	//connexion
	bbs.Client = bitbucket.NewOAuthClientCredentials(bbs.key, bbs.secret)

	//Set url
	if err := bbs.bitbucket_set_url(server); err != nil{
		ret.Errorf("Invalid url. %s", err)
		return nil
	}

	//verif
    userProfil, err := bbs.Client.User.Profile()
	_ = userProfil
	if err != nil {
		ret.Errorf("Unable to get owner of given token. %s", err)
		return nil
	} else {
		ret.StatusAdd("Connection successful.")
	}
	return bbs.Client
}*/

/*func (bbs *BitbucketPlugin) bitbucket_set_url(server string) (err error) {
	//gl_url := ""
	if bbs.bitbucket_source.Urls == nil {
		bbs.bitbucket_source.Urls = make(map[string]string)
	}
	//...
	return //TODO
}*/

/*func (req *CreateReq) InitTeam(bbs *BitbucketPlugin) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		bbs.SetTeam(app)
		ret = true
	}
	return
}

func (bbs *BitbucketPlugin) SetTeam(fromApp AppInstanceStruct) {
	if team := fromApp.Team; team == ""{
		bbs.bitbucketDeploy.Team =fromApp.ForjjTeam
	} else {
		bbs.bitbucketDeploy.Team = team
	}
	if team := fromApp.ProductionTeam; team == ""{
		bbs.bitbucketDeploy.ProdTeam = fromApp.ForjjTeam
	} else {
		bbs.bitbucketDeploy.ProdTeam = team
	}
	bbs.bitbucket_source.ProdTeam = bbs.bitbucketDeploy.ProdTeam
}*/

/*func (bbs *BitbucketPlugin) create_yaml_data(req *CreateReq, ret *goforjj.PluginData) error{
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
}*/


func (r *RepoInstanceStruct) IsValid(repo_name string, ret *goforjj.PluginData) (valid bool){ // /!\ change struct (project)
	if r.Name == "" {
		ret.Errorf("Invalid repo '%s'. Name is empty.", repo_name)
		return
	}
	if r.Name != repo_name {
		ret.Errorf("Invalid repo '%s'. Name must be equal to '%s'. But the repo name is set to '%s'.", repo_name, repo_name, r.Name)
		return
	}
	valid = true
	return
}

/*func (bbs *BitbucketPlugin) SetRepo(repo *RepoInstanceStruct, isInfra, isDeployable bool) { //SetRepo
	//upstream := gls.DefineRepoUrls(repo.Name)

	owner := gls.gitlabDeploy.Group
	if isInfra {
		owner = gls.gitlabDeploy.ProdGroup
	}

	//set it, found or not
	r := RepositoryStruct{}
	r.set(repo)
	bbs.bitbucketDeploy.Repos[repo.Name] = r
}*/

func (r *RepositoryStruct) set(repo *RepoInstanceStruct) *RepositoryStruct{
	if r == nil {
		r = new(RepositoryStruct)
	}
	r.Name = repo.Name
	return r
}

/*func (bbs *BitbucketPlugin) repos_exists(ret *goforjj.PluginData) (err error) {
	clientRepos := bbs.Client.Repositories //Repos
	//client, err := bbs.Client.User.Profile() //Get current user profile
	
	//loop
	for name, repo_data := range bbs.bitbucketDeploy.Repos{

		//
		RepoOptions := &bitbucket.RepositoryOptions{
	    	Owner: repo_data.Owner,
	    	RepoSlug: name,
	    }

		//on voit si on trouve le repo X sur gitlab
		if found_project, e := clientRepos.Repository.Get(RepoOptions); e == nil{
			//Si trouvé: err
			if err == nil && name == bbs.app.ForjjInfra {
				err = fmt.Errorf("Infra repo '%s' already exist in bitbucket server.", name)
			}
			repo_data.exist = true
			if repo_data.remotes == nil{
				repo_data.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				repo_data.branchConnect = make(map[string]string)
			}

			//Get ssh and https remte
			ssh := "";
			url := "";
			remotes := found_project.Links["clone"].([]interface{})
			if remotes[0].(map[string]interface{})["name"].(string) == "https"{
				url = remotes[0].(map[string]interface{})["href"].(string)
				ssh = remotes[1].(map[string]interface{})["href"].(string)
			} else {
				url = remotes[1].(map[string]interface{})["href"].(string)
				ssh = remotes[0].(map[string]interface{})["href"].(string)
			}

			repo_data.remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: ssh,
				Url: url,
			}
			repo_data.branchConnect["master"] = "origin/master"
		}
		ret.Repos[name] = goforjj.PluginRepo{ //Project ?
			Name: 		repo_data.Name,
			Exist: 		repo_data.exist,
			Remotes: 		repo_data.remotes,
			BranchConnect: 		repo_data.branchConnect,
			//Owner: 		repo_data.Owner,
		}

	}

	return
}*/

/*func (bbs *BitbucketPlugin) checkSourcesExistence(when string) (err error){
	log.Print("Checking Infrastructure code existence.")
	sourceRepo := bbs.source_path
	sourcePath := path.Join(sourceRepo, bbs.instance)
	bbs.sourceFile = path.Join(sourcePath, bitbucket_file)

	deployRepo := path.Join(bbs.deployMount, bbs.deployTo)
	deployBase := path.Join(deployRepo, bbs.instance)

	bbs.deployFile = path.Join(deployBase, bitbucket_file)

	bbs.gitFile = path.Join(bbs.instance, bitbucket_file)

	switch when {
		case "create":
			if _, err := os.Stat(sourcePath); err != nil{
				if err = os.MkdirAll(sourcePath, 0755); err != nil{
					return fmt.Errorf("Unable to create '%s'. %s", sourcePath, err)
				}
			}

			if _, err := os.Stat(deployRepo); err != nil{
				return fmt.Errorf("Unable to create '%s'. Forjj must create it. %s", deployRepo, err)
			}

			if _, err := os.Stat(bbs.sourceFile); err == nil{
				return fmt.Errorf("Unable to create the gitlab configuration which already exist.\nUse 'update' to update it "+"(or update %s), and 'maintain' to update your github service according to his configuration.",path.Join(bbs.instance, bitbucket_file))
			}

			if _, err := os.Stat(deployBase); err != nil{
				if err = os.Mkdir(deployBase, 0755); err != nil{
					return fmt.Errorf("Unable to create '%s'. %s", deployBase, err)
				}
			}
			return
		
		case "update":
			log.Printf("TODO UPDATE")
	}
	return
}*/

/*func (bbs *BitbucketPlugin) save_yaml(in interface{}, file string) (Updated bool, _ error){
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
}*/

// DoCreate --> plugin tasks for forjj create
// req_data contains the request data posted by forjj. Structure generated from 'bitbucket.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {
	
	//Get instance name
	instance := req.Forj.ForjjInstanceName

	//init BitbucketPlugin
	bbs := BitbucketPlugin {
		source_path: 	req.Forj.ForjjSourceMount,
		deployMount: 	req.Forj.ForjjDeployMount,
		instance: 		req.Forj.ForjjInstanceName,
		deployTo: 		req.Forj.ForjjDeploymentEnv,
		key: 			req.Objects.App[instance].Key,
		secret: 		req.Objects.App[instance].Secret,
		team:			req.Objects.App[instance].Team,
	}

	log.Printf("Checking parameters : %#v", bbs)

	//Check key, secret and source
	check := make(map[string]bool)
	check["key"] = true
	check["secret"] = true
	check["source"] = true

	//ensure source path is writeable && key/secret isn't empty
	if bbs.verify_req_fails(ret, check){
		return
	}

	//verify if instance existe
	if a, found := req.Objects.App[instance]; !found{
		ret.Errorf("Internal issue. Forjj has not given the Application information for '%s'. Aborted.")
		return
	} else {
		bbs.app = &a
	}

	//Check Gitlab connection
	log.Println("Checking bitbucket connection.")
	ret.StatusAdd("Connect to bitbucket...")

	if git := bbs.bitbucket_connect("X", ret); git == nil{ //!\\
		return
	}

	// Init Group of project
	if !req.InitTeam(&bbs){
		ret.Errorf("Internal Error. Unable to define the team.")
	}

	//Create yaml data for maintain function
	if err := bbs.create_yaml_data(req, ret); err != nil {
		ret.Errorf("Unable to create. %s",err)
		return
	}

	//repos exist ?
	if err := bbs.repos_exists(ret); err != nil {
		ret.Errorf("%s\nUnable to 'create' your forge when gitlab already has an infra project created. Clone it and use 'update' instead.", err)
		return 419
	}

	//CheckSourceExistence
	if err := bbs.checkSourcesExistence("create"); err != nil{
		ret.Errorf("%s\nUnable to 'create' your forge", err)
		return
	}
	ret.StatusAdd("Environment checked. Ready to be created.")

	//Path in ctxt git
	gitFile := path.Join(bbs.instance, bitbucket_file)

	//Save bitbucket source
	if _, err := bbs.save_yaml(&bbs.bitbucket_source, bbs.sourceFile); err != nil {
		ret.Errorf("%s", err)
		return
	}
	log.Printf(ret.StatusAdd("Configuration saved in source repo '%s' (%s).", gitFile, bbs.source_path))

	//Save bitbucket deploy
	if _, err := bbs.save_yaml(&bbs.bitbucketDeploy, bbs.deployFile); err != nil{
		ret.Errorf("%s", err)
		return
	}
	log.Printf(ret.StatusAdd("Configuration saved in deploy repo '%s' (%s).", gitFile, path.Join(bbs.deployMount, bbs.deployTo)))

	//Build final post answer
	for k, v := range bbs.bitbucket_source.Urls{
		ret.Services.Urls[k] = v
	}

	//API by forjj
	ret.Services.Urls["api_url"] = bbs.bitbucket_source.Urls["bitbucket-base-url"]

	ret.CommitMessage = fmt.Sprint("Bitbucket configuration created.")

	ret.AddFile(goforjj.FilesSource, gitFile)
	ret.AddFile(goforjj.FilesDeploy, gitFile)


	ret.StatusAdd("end");
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
	// This is where you shoud write your Maintain code. Following lines are typical code for a basic plugin.
	//TODO
	return
}
