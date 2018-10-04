package main

import(
	"os"
	"fmt"
	"path"
	"log"

	"github.com/forj-oss/goforjj"
	"github.com/ktrysmt/go-bitbucket"
)

func (bbs *BitbucketPlugin) bitbucket_connect(server string, ret *goforjj.PluginData) *bitbucket.Client{
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
}

func (req *CreateReq) InitTeam(bbs *BitbucketPlugin) (ret bool) {
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
}

func (bbs *BitbucketPlugin) ensureTeamExists(ret *goforjj.PluginData) (s bool){
	//TODO
	return																   
}

func IsNewForge() {
	//TODO
}

func (bbs *BitbucketPlugin) bitbucket_set_url(server string) (err error) {
	//gl_url := ""
	if bbs.bitbucket_source.Urls == nil {
		bbs.bitbucket_source.Urls = make(map[string]string)
	}
	//...
	return //TODO
}

func ensureExists () {
	//TODO
}

func (bbs *BitbucketPlugin) repos_exists(ret *goforjj.PluginData) (err error) {
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
}

func (bbs *BitbucketPlugin) checkSourcesExistence(when string) (err error){
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
}
