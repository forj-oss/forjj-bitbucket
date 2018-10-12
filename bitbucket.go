package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/forj-oss/goforjj"
	"github.com/ktrysmt/go-bitbucket"
)

func (bbs *BitbucketPlugin) bitbucketConnect(server string, ret *goforjj.PluginData) *bitbucket.Client {
	//connexion
	bbs.Client = bitbucket.NewOAuthClientCredentials(bbs.key, bbs.secret)

	//Set url
	if err := bbs.bitbucketSetUrl(server); err != nil {
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
}

//InitTeam ...
func (req *CreateReq) InitTeam(bbs *BitbucketPlugin) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found {
		bbs.SetTeam(app)
		ret = true
	}
	return
}

//SetTeam ...
func (bbs *BitbucketPlugin) SetTeam(fromApp AppInstanceStruct) {
	if team := fromApp.Team; team == "" {
		bbs.bitbucketDeploy.Team = fromApp.ForjjTeam
	} else {
		bbs.bitbucketDeploy.Team = team
	}
	if team := fromApp.ProductionTeam; team == "" {
		bbs.bitbucketDeploy.ProdTeam = fromApp.ForjjTeam
	} else {
		bbs.bitbucketDeploy.ProdTeam = team
	}
	bbs.bitbucketSource.ProdTeam = bbs.bitbucketDeploy.ProdTeam
}

//ensureTeamExists
func (bbs *BitbucketPlugin) ensureTeamExists(ret *goforjj.PluginData) (s bool) {
	//TODO
	return
}

//IsNewForge TODO
func (bbs *BitbucketPlugin) IsNewForge(ret *goforjj.PluginData) (_ bool) {
	c := bbs.Client.Repositories

	//loop on list of repos, and ensure they exist with minimal config and rights
	for _, repo := range bbs.bitbucketDeploy.Repos {
		if !repo.Infra {
			continue
		}

		//Get username
		userProfil, _ := bbs.Client.User.Profile() //!\\
		jsonMap := userProfil.(map[string]interface{})
		//Repository Options
		ro := &bitbucket.RepositoryOptions{
			Owner:    jsonMap["username"].(string),
			RepoSlug: "name",
		}

		if _, e := c.Repository.Get(ro); e == nil {
			ret.Errorf("Unable to identify the infra repository. Unknown issue: %s", e)
		} else {
			//bbs.newForge = (resp.StatusCode != 200) TODO
		}
		return true
	}

	ret.Errorf("Unable to identify the infra repository. At least, one repo must be identified with"+"`%s` in %s. You can use Forjj update to fix this.", "Infra: true", "bitbucket")
	return
}

//bitbucketSetUrl TODO
func (bbs *BitbucketPlugin) bitbucketSetUrl(server string) (err error) {
	bbUrl := ""

	if bbs.bitbucketSource.Urls == nil {
		bbs.bitbucketSource.Urls = make(map[string]string)
	}

	if !bbs.maintainCtxt {
		if server == "" { // || ?
			bbs.bitbucketSource.Urls["bitbucket-base-url"] = "https://bitbucket.com/"
			bbs.bitbucketSource.Urls["bitbucket-url"] = "https://bitbucket.com"
			bbs.bitbucketSource.Urls["bitbucket-ssh"] = "git@bitbucket.com:"
		} else {
			//set from serveur // ! \\ TODO
			server = "bitbucket.com"
			bbUrl = "https://" + server + "/api/v4/"
			bbs.bitbucketSource.Urls["bitbucket-url"] = "https://bitbucket.com"
			bbs.bitbucketSource.Urls["bitbucket-ssh"] = "git@bitbucket.com:"
		}
	} else {
		//maintain context
		bbs.bitbucketSource.Urls = bbs.bitbucketDeploy.Urls
		bbUrl = bbs.bitbucketSource.Urls["bitbucket-base-url"]
	}

	if bbUrl == "" {
		return
	}

	//err = bbs.Client.SetBaseURL(bbUrl) TODO SETBASEURL

	/*if err != nil{
		return
	}*/

	return
}

//ensureExists TODO
func (r *RepositoryStruct) ensureExists(bbs *BitbucketPlugin, ret *goforjj.PluginData) error {
	//test existence
	clientRepos := bbs.Client.Repositories
	userProfil, err := bbs.Client.User.Profile()
	if err != nil {
		ret.Errorf("Unable to identify owner.")
		return err
	}
	user := userProfil.(map[string]interface{})

	RepoOptions := &bitbucket.RepositoryOptions{
		Owner:    user["username"].(string),
		RepoSlug: r.Name,
	}

	_, e := clientRepos.Repository.Get(RepoOptions)

	if e != nil {
		//Create
		_, er := clientRepos.Repository.Create(RepoOptions)
		if er != nil {
			ret.Errorf("Unable to create '%s'. %s.", r.Name, er)
			return er
		}
		log.Printf(ret.StatusAdd("Repo '%s': created", r.Name))
	} else {
		//Update TODO
	}

	//...

	return nil
}

//reposExists TODO
func (bbs *BitbucketPlugin) reposExists(ret *goforjj.PluginData) (err error) {
	clientRepos := bbs.Client.Repositories //Repos
	//client, err := bbs.Client.User.Profile() //Get current user profile

	//loop
	for name, repoData := range bbs.bitbucketDeploy.Repos {

		//
		RepoOptions := &bitbucket.RepositoryOptions{
			Owner:    repoData.Owner,
			RepoSlug: name,
		}

		//on voit si on trouve le repo X sur bitbucket
		if foundProject, e := clientRepos.Repository.Get(RepoOptions); e == nil {
			//Si trouvé: err
			if err == nil && name == bbs.app.ForjjInfra {
				err = fmt.Errorf("Infra repo '%s' already exist in bitbucket server.", name)
			}
			repoData.exist = true
			if repoData.remotes == nil {
				repoData.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				repoData.branchConnect = make(map[string]string)
			}

			//Get ssh and https remte
			ssh := ""
			url := ""
			remotes := foundProject.Links["clone"].([]interface{})
			if remotes[0].(map[string]interface{})["name"].(string) == "https" {
				url = remotes[0].(map[string]interface{})["href"].(string)
				ssh = remotes[1].(map[string]interface{})["href"].(string)
			} else {
				url = remotes[1].(map[string]interface{})["href"].(string)
				ssh = remotes[0].(map[string]interface{})["href"].(string)
			}

			repoData.remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: ssh,
				Url: url,
			}
			repoData.branchConnect["master"] = "origin/master"
		}
		ret.Repos[name] = goforjj.PluginRepo{ //Project ?
			Name:          repoData.Name,
			Exist:         repoData.exist,
			Remotes:       repoData.remotes,
			BranchConnect: repoData.branchConnect,
			//Owner: 		repoData.Owner,
		}

	}

	return
}

func (bbs *BitbucketPlugin) checkSourcesExistence(when string) (err error) {
	log.Print("Checking Infrastructure code existence.")
	sourceRepo := bbs.sourcePath
	sourcePath := path.Join(sourceRepo, bbs.instance)
	bbs.sourceFile = path.Join(sourcePath, bitbucketFile)

	deployRepo := path.Join(bbs.deployMount, bbs.deployTo)
	deployBase := path.Join(deployRepo, bbs.instance)

	bbs.deployFile = path.Join(deployBase, bitbucketFile)

	bbs.gitFile = path.Join(bbs.instance, bitbucketFile)

	switch when {
	case "create":
		if _, err := os.Stat(sourcePath); err != nil {
			if err = os.MkdirAll(sourcePath, 0755); err != nil {
				return fmt.Errorf("Unable to create '%s'. %s", sourcePath, err)
			}
		}

		if _, err := os.Stat(deployRepo); err != nil {
			return fmt.Errorf("Unable to create '%s'. Forjj must create it. %s", deployRepo, err)
		}

		if _, err := os.Stat(bbs.sourceFile); err == nil {
			return fmt.Errorf("Unable to create the bitbucket configuration which already exist.\nUse 'update' to update it "+"(or update %s), and 'maintain' to update your github service according to his configuration.", path.Join(bbs.instance, bitbucketFile))
		}

		if _, err := os.Stat(deployBase); err != nil {
			if err = os.Mkdir(deployBase, 0755); err != nil {
				return fmt.Errorf("Unable to create '%s'. %s", deployBase, err)
			}
		}
		return

	case "update":
		log.Printf("TODO UPDATE")
	}
	return
}

//setTeamHooks TODO
func (bbs *BitbucketPlugin) setTeamHooks(teamHookDisabled, repoHookDisabled, whPolicy string, hooks map[string]WebhooksInstanceStruct) {
	//
	if b, err := strconv.ParseBool(teamHookDisabled); err != nil {
		log.Printf("Team webhook disabled: invalid boolean: %s", teamHookDisabled)
		bbs.bitbucketDeploy.NoTeamHook = true
	} else {
		bbs.bitbucketDeploy.NoTeamHook = b
	}
	if bbs.bitbucketDeploy.WebHooks == nil {
		bbs.bitbucketDeploy.WebHooks = make(map[string]WebHookStruct)
	}

	if b, err := strconv.ParseBool(repoHookDisabled); err != nil {
		log.Printf("Team webhook disabled: invalid boolean: %s", repoHookDisabled)
	} else {
		bbs.bitbucketDeploy.NoRepoHook = b
	}

	if v := inStringList(whPolicy, "manage", "sync"); v == "" || v == "sync" {
		if whPolicy != "" {
			log.Printf("Invalid value '%s' for 'WebhooksManagement'. Set it to 'sync'.", whPolicy)
		} else {
			log.Print("WebhooksManagement is set by default to 'sync'.")
		}
		bbs.bitbucketDeploy.WebHookPolicy = ""
	} else {
		bbs.bitbucketDeploy.WebHookPolicy = v
	}

	if bbs.bitbucketDeploy.NoTeamHook {
		return
	}

	for name, hook := range hooks {
		if hook.Team == "false" {
			continue
		}
		data := WebHookStruct{
			Url:     hook.Url,
			Events:  strings.Split(hook.Events, ","),
			Enabled: hook.Enabled,
		}
		if v, err := strconv.ParseBool(hook.SslCheck); err == nil {
			data.SSLCheck = v
			log.Printf("SSL Check '%s' => %t", name, v)
		} else {
			log.Printf("SSLCheck has an invalid boolean string representation '%s'. Ignore. SSL Check is set to true.", name)
			data.SSLCheck = true
		}
		bbs.bitbucketDeploy.WebHooks[name] = data
	}
	if len(bbs.bitbucketDeploy.WebHooks) > 0 && bbs.bitbucketDeploy.WebHookPolicy == "sync" {
		bbs.bitbucketDeploy.WebHookPolicy = ""
	}
}
