package conf

import "encoding/json"

type RepoData struct {
	Repo       Repo       `json:"Repo"`
	RepoSource RepoSource `json:"RepoSource"`
}
type Repo struct {
	RepoNamespace string `json:"RepoNamespace"`
	RepoName      string `json:"RepoName"`
	Summary       string `json:"Summary"`
	RepoType      string `json:"RepoType"`
	RepoBuildType string `json:"RepoBuildType"`
}
type Source struct {
	SourceRepoType      string `json:"SourceRepoType"`
	SourceRepoNamespace string `json:"SourceRepoNamespace"`
	SourceRepoName      string `json:"SourceRepoName"`
}
type BuildConfig struct {
	IsAutoBuild    bool `json:"IsAutoBuild"`
	IsOversea      bool `json:"IsOversea"`
	IsDisableCache bool `json:"IsDisableCache"`
}
type RepoSource struct {
	Source      Source      `json:"Source"`
	BuildConfig BuildConfig `json:"BuildConfig"`
}

// 创建阿里云仓库所需 request 参数
func GenRepoData(data RepoData) string {
	res, _ := json.Marshal(data)
	return string(res)
}