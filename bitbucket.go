package main

import(
	"os"
	"fmt"
	"path"
	"log"

	"github.com/forj-oss/goforjj"
	"github.com/ktrysmt/go-bitbucket"
)

func (bbs *BitbucketPlugin) bitbucketConnect(server string, ret *goforjj.PluginData) *bitbucket.Client{
	//connexion
	bbs.Client = bitbucket.NewOAuthClientCredentials(bbs.key, bbs.secret)

	//Set url
	if err := bbs.bitbucketSetUrl(server); err != nil{
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
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		bbs.SetTeam(app)
		ret = true
	}
	return
}

//SetTeam ...
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
	bbs.bitbucketSource.ProdTeam = bbs.bitbucketDeploy.ProdTeam
}

//ensureTeamExists 
func (bbs *BitbucketPlugin) ensureTeamExists(ret *goforjj.PluginData) (s bool){
	//TODO
	return																   
}

//IsNewForge TODO
func (bbs *BitbucketPlugin) IsNewForge(ret *goforjj.PluginData) (_ bool){
	c := bbs.Client.Repositories

	//loop on list of repos, and ensure they exist with minimal config and rights
	for name, repo := range bbs.bitbucketDeploy.Repos{
		if !repo.Infra {
			continue
		}

		//Get username
		userProfil, _ := bbs.Client.User.Profile() //!\\
		jsonMap := userProfil.(map[string]interface{})
		//Repository Options
		ro := &bitbucket.RepositoryOptions{
			Owner: jsonMap["username"].(string),
			RepoSlug: name,
		}
		if resp, e := c.Repository.Get(ro); e != nil && resp == nil {
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

	if bbUrl == ""{
		return
	}

	//err = bbs.Client.SetBaseURL(bbUrl) TODO SETBASEURL
	
	if err != nil{
		return
	}

	return
}

//ensureExists TODO
func ensureExists () {
	//TODO
}

//reposExists TODO
func (bbs *BitbucketPlugin) reposExists(ret *goforjj.PluginData) (err error) {
	clientRepos := bbs.Client.Repositories //Repos
	//client, err := bbs.Client.User.Profile() //Get current user profile
	
	//loop
	for name, repoData := range bbs.bitbucketDeploy.Repos{

		//
		RepoOptions := &bitbucket.RepositoryOptions{
	    	Owner: repoData.Owner,
	    	RepoSlug: name,
	    }

		//on voit si on trouve le repo X sur bitbucket
		if foundProject, e := clientRepos.Repository.Get(RepoOptions); e == nil{
			//Si trouvé: err
			if err == nil && name == bbs.app.ForjjInfra {
				err = fmt.Errorf("Infra repo '%s' already exist in bitbucket server.", name)
			}
			repoData.exist = true
			if repoData.remotes == nil{
				repoData.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				repoData.branchConnect = make(map[string]string)
			}

			//Get ssh and https remte
			ssh := "";
			url := "";
			remotes := foundProject.Links["clone"].([]interface{})
			if remotes[0].(map[string]interface{})["name"].(string) == "https"{
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
			Name: 		repoData.Name,
			Exist: 		repoData.exist,
			Remotes: 		repoData.remotes,
			BranchConnect: 		repoData.branchConnect,
			//Owner: 		repoData.Owner,
		}

	}

	return
}

func (bbs *BitbucketPlugin) checkSourcesExistence(when string) (err error){
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
			if _, err := os.Stat(sourcePath); err != nil{
				if err = os.MkdirAll(sourcePath, 0755); err != nil{
					return fmt.Errorf("Unable to create '%s'. %s", sourcePath, err)
				}
			}

			if _, err := os.Stat(deployRepo); err != nil{
				return fmt.Errorf("Unable to create '%s'. Forjj must create it. %s", deployRepo, err)
			}

			if _, err := os.Stat(bbs.sourceFile); err == nil{
				return fmt.Errorf("Unable to create the bitbucket configuration which already exist.\nUse 'update' to update it "+"(or update %s), and 'maintain' to update your github service according to his configuration.",path.Join(bbs.instance, bitbucketFile))
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
}
