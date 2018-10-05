package main

import(
	"github.com/forj-oss/goforjj"
)

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

	//maintain
	Infra 			bool 				`yaml:",omitempty"`
	Role 			string 				`yaml:",omitempty"`
	IsDeployable 		bool

}

//IsValid ...
func (r *RepoInstanceStruct) IsValid(repoName string, ret *goforjj.PluginData) (valid bool){
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

//set
func (r *RepositoryStruct) set(repo *RepoInstanceStruct, remotes map[string]goforjj.PluginRepoRemoteUrl, branchConnect map[string]string, isInfra, IsDeployable bool, owner string) *RepositoryStruct{
	if r == nil {
		r = new(RepositoryStruct)
	}
	r.Name = repo.Name
	return r

	//TODO
}
