package githubApiSelf

import (
	"context"
	"fmt"
	"github.com/hongqchen/imagectl/utils/conf"
	"github.com/hongqchen/imagectl/utils/log"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"strings"
)

var (
	GithubClient *github.Client
)

type GithubApi struct {
	RepoName string
	Err      error
}

// Connected 提供 access token，建立 github 连接，
func (ga *GithubApi) connected(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

// 在 github 创建镜像仓库，并返回仓库所属用户名
func (ga *GithubApi) CreateRepository(ctx context.Context, repoName string) {
	_, _, err := GithubClient.Repositories.Create(ctx, "",
		&github.Repository{
			Name: &repoName,
		})
	if err != nil {
		ga.Err = err
		return
	}

	// 记录仓库名
	ga.RepoName = repoName
}

// 在 github 仓库创建 Dockerfile 文件
func (ga *GithubApi) CreateFile(ctx context.Context, content string) {
	// 判断前序步骤是否出错，防止缺失信息，导致任务失败
	if ga.Err != nil {
		return
	}
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("创建 Dockerfile 文件"),
		Content: []byte(content),
	}
	_, _, err := GithubClient.Repositories.CreateFile(ctx, conf.Global.Github.Namespace, ga.RepoName, "Dockerfile", opts)
	if err != nil {
		ga.Err = err
		return
	}
}

// 判断仓库是否存在
func IsExistsRepo(ctx context.Context, owner, repoName string) (bool, error) {
	log.Logger.Debugf("[Github] 判断镜像仓库 %s 是否存在", repoName)
	_, _, err := GithubClient.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 删除 github 仓库
func DeleteRepo(ctx context.Context, owner, repoName string) error {
	isexists, err := IsExistsRepo(ctx, owner, repoName)
	if err != nil {
		return err
	}
	if isexists {
		_, err = GithubClient.Repositories.Delete(ctx, owner, repoName)
		if err != nil {
			return err
		}
	}

	return nil
}

// 建立连接
func GetGithubConnection() {
	githubObj := new(GithubApi)
	GithubClient = githubObj.connected(conf.Global.Github.Access_token)
}

// 创建仓库和 Dockerfile
func CallGithubApi(ctx context.Context, image map[string]string) {
	log.Logger.Debugf("[Github] 镜像仓库 %s:%s 创建中", image["image_name"], image["image_version"])
	githubObj := new(GithubApi)
	githubObj.CreateRepository(ctx, image["image_name"])
	githubObj.CreateFile(ctx, fmt.Sprintf("FROM %s", image["image_url"]))
	if githubObj.Err != nil {
		if strings.Contains(githubObj.Err.Error(), "401 Bad credentials") {
			log.Logger.Warn("[Github] secret 验证失败，请检查 token 是否正确!")
			return
		}
		log.Logger.Error(githubObj.Err.Error())
		return
	}
}
