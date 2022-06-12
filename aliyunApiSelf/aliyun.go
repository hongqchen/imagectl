package aliyunApiSelf

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	cr20160607 "github.com/hongqchen/alisdk/client"
	"github.com/hongqchen/imagectl/githubApiSelf"
	"github.com/hongqchen/imagectl/utils/conf"
	"github.com/hongqchen/imagectl/utils/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var (
	AliyunClient *cr20160607.Client
)

type AliyunApi struct {
	RepoName    string
	RepoVersion string
}

// 连接至 aliyun，创建 api client
func (aa *AliyunApi) createClient(accessKeyId, accessKeySecret *string, region string) error {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String(region)

	result, err := cr20160607.NewClient(config)
	if err != nil {
		return err
	}
	AliyunClient = result
	return nil
}

// 查询镜像是否存在
func (aa *AliyunApi) getRepoIsExist() (bool, error) {
	// 判断 namespace 是否存在
	_, err := AliyunClient.GetNamespace(&conf.Global.Aliyun.Namespace)
	if err != nil {
		if sdkErr, isok := err.(*tea.SDKError); isok {
			if strings.Contains(tea.StringValue(sdkErr.Code), "NAMESPACE_NOT_EXIST") {
				return false, errors.New("namespace 不存在!")
			} else if strings.Contains(tea.StringValue(sdkErr.Code), "InvalidAction.Mismatch") {
				return false, errors.New("namespace 格式错误!")
			}
		}

		return false, err
	}

	_, err = AliyunClient.GetRepo(&conf.Global.Aliyun.Namespace, &aa.RepoName)
	if err != nil {
		if sdkErr, isok := err.(*tea.SDKError); isok {
			if strings.Contains(tea.StringValue(sdkErr.Code), "AccessKeyId") ||
				strings.Contains(tea.StringValue(sdkErr.Code), "SignatureDoesNotMatch") {
				return false, errors.New("access key 格式错误!")
			} else if strings.Contains(tea.StringValue(sdkErr.Code), "REPO_NOT_EXIST") {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

// 创建镜像仓库
func (aa *AliyunApi) createRepository(ctx context.Context, image map[string]string) error {
	// 镜像仓库已经存在，跳过创建过程
	log.Logger.Debugf("[阿里云] 判断镜像仓库 %s:%s 是否存在", aa.RepoName, aa.RepoVersion)
	isExist, err := aa.getRepoIsExist()
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}

	log.Logger.Debugf("[阿里云] 镜像仓库 %s:%s 创建中", image["image_name"], image["image_version"])
	isexists, err := githubApiSelf.IsExistsRepo(ctx, conf.Global.Github.Namespace, image["image_name"])
	if err != nil {
		return err
	}
	if !isexists {
		githubApiSelf.CallGithubApi(ctx, image)
	}

	repodata := conf.GenRepoData(conf.RepoData{
		Repo: conf.Repo{
			RepoNamespace: conf.Global.Aliyun.Namespace,
			RepoName:      aa.RepoName,
			Summary:       fmt.Sprintf("%s 拉取", aa.RepoName),
			RepoType:      conf.Global.Aliyun.Repo_type,
			RepoBuildType: "AUTO_BUILD",
		},
		RepoSource: conf.RepoSource{
			Source: conf.Source{
				SourceRepoType:      "GITHUB",
				SourceRepoNamespace: conf.Global.Github.Namespace,
				SourceRepoName:      aa.RepoName,
			},
			BuildConfig: conf.BuildConfig{
				IsAutoBuild:    false,
				IsOversea:      true,
				IsDisableCache: false,
			},
		},
	})
	_, err = AliyunClient.CreateRepo(repodata)
	if err != nil {
		return err
	}
	return nil
}

// 查询镜像构建规则是否已经存在，以版本信息为主
func (aa *AliyunApi) getRepoRuleIsExist() (int, bool, error) {
	res, err := AliyunClient.GetRepoBuildRuleList(&conf.Global.Aliyun.Namespace, &aa.RepoName)
	if err != nil {
		return 0, false, err
	}
	for _, value := range *res.Body.Data.BuildRules {
		if *value.ImageTag == aa.RepoVersion {
			return *value.BuildRuleID, true, nil
		}
	}
	return 0, false, nil
}

// 创建镜像构建规则
func (aa *AliyunApi) createBuildRule() (bool, string, error) {
	log.Logger.Debugf("[阿里云] 判断镜像 %s:%s 是否存在", aa.RepoName, aa.RepoVersion)
	buildStatus, err := aa.getAllBuildStatus()
	if err != nil {
		return false, "", err
	}
	if buildStatus == "SUCCESS" {
		log.Logger.Debugf("[阿里云] 镜像 %s:%s 已存在，直接返回镜像地址", aa.RepoName, aa.RepoVersion)
		return true, fmt.Sprintf("registry.cn-%s.aliyuncs.com/%s/%s:%s\n", conf.Global.Aliyun.Region,
			conf.Global.Aliyun.Namespace, aa.RepoName, aa.RepoVersion), nil
	}

	// 判断当前版本构建规则是否存在，存在则跳过创建过程
	log.Logger.Debugf("[阿里云] 查询镜像 %s:%s 构建规则是否存在", aa.RepoName, aa.RepoVersion)
	buildRuleId, isExist, err := aa.getRepoRuleIsExist()
	if err != nil {
		return false, "", err
	}
	if isExist {
		log.Logger.Debugf("[阿里云] 镜像 %s:%s 构建规则已存在，忽略创建请求", aa.RepoName, aa.RepoVersion)
		return false, strconv.Itoa(buildRuleId), nil
	}

	log.Logger.Debugf("[阿里云] 镜像 %s:%s 构建规则不存在，规则创建中", aa.RepoName, aa.RepoVersion)
	ruleData := conf.GenRuleData(conf.RuleData{
		BuildRule: conf.BuildRule{
			DockerfileLocation: "/",
			DockerfileName:     "Dockerfile",
			PushType:           "GIT_BRANCH",
			PushName:           "main",
			ImageTag:           aa.RepoVersion,
			Tag:                aa.RepoVersion,
		},
	})

	res, err := AliyunClient.CreateRepoBuildRule(&conf.Global.Aliyun.Namespace, &aa.RepoName, ruleData)
	if err != nil {
		if SDKErrorInfo, isok := err.(*tea.SDKError); isok {
			return false, "", errors.New(tea.StringValue(SDKErrorInfo.Message))
		}

		return false, "", err
	}
	return false, tea.StringValue(res.Body.Data.BuildRuleId), nil
}

// 开始构建
func (aa *AliyunApi) startBuild(ctx context.Context, buildRuleId string, image map[string]string) (string, error) {
	isexists, err := githubApiSelf.IsExistsRepo(ctx, conf.Global.Github.Namespace, image["image_name"])
	if err != nil {
		return "", err
	}
	if !isexists {
		githubApiSelf.CallGithubApi(ctx, image)
	}

	_, err = AliyunClient.StartRepoBuildByRule(&conf.Global.Aliyun.Namespace, &aa.RepoName, &buildRuleId)
	if err != nil {
		return "", err
	}

	log.Logger.Debugf("[阿里云] 镜像 %s:%s 构建状态查询", aa.RepoName, aa.RepoVersion)
	imageURL, err := aa.getBuildStatusLoop()
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

// 查询最新发起的构建任务状态
func (aa *AliyunApi) getBuildStatus() (string, error) {
	buildListRequest := cr20160607.GetRepoBuildListRequest{
		Page:     tea.Int32(1),
		PageSize: tea.Int32(100),
	}

	res, err := AliyunClient.GetRepoBuildList(&conf.Global.Aliyun.Namespace, &aa.RepoName, &buildListRequest)
	if err != nil {
		return "", err
	}
	resSlice := *res.Body.Data.Builds
	return *resSlice[0].BuildStatus, nil
}

// 查看以往构建记录任务状态
func (aa *AliyunApi) getAllBuildStatus() (string, error) {
	buildListRequest := cr20160607.GetRepoBuildListRequest{
		Page:     tea.Int32(1),
		PageSize: tea.Int32(100),
	}
	res, err := AliyunClient.GetRepoBuildList(&conf.Global.Aliyun.Namespace, &aa.RepoName, &buildListRequest)
	if err != nil {
		return "", err
	}

	tmpBuildItems := []cr20160607.Builds{}
	for _, buildList := range *res.Body.Data.Builds {
		// 匹配当前版本的构建记录
		if *buildList.Image.Tag != aa.RepoVersion {
			continue
		}
		tmpBuildItems = append(tmpBuildItems, buildList)
	}

	for _, value := range tmpBuildItems {
		if *value.BuildStatus == "SUCCESS" {
			return "SUCCESS", nil
		}
	}
	return "FAILED", nil
}

func (aa *AliyunApi) getBuildStatusLoop() (string, error) {
	buildStatus := ""
	for {
		buildStatusTmp, err := aa.getBuildStatus()
		if err != nil {
			return "", err
		}
		buildStatus = buildStatusTmp

		if buildStatusTmp == "SUCCESS" || buildStatusTmp == "FAILED" {
			break
		}
		log.Logger.Debugf("[阿里云] %s:%s 构建状态: %s，等待 5s 后再次查询任务状态", aa.RepoName, aa.RepoVersion, buildStatusTmp)
		time.Sleep(5 * time.Second)
	}

	if buildStatus == "SUCCESS" {
		return fmt.Sprintf("registry.cn-%s.aliyuncs.com/%s/%s:%s\n", conf.Global.Aliyun.Region,
			conf.Global.Aliyun.Namespace, aa.RepoName, aa.RepoVersion), nil
	}
	errMessage := fmt.Sprintf("错误原因登录 https://cr.console.aliyun.com/repository/%s/%s/%s/build 查看", conf.Global.Aliyun.Region,
		conf.Global.Aliyun.Namespace, aa.RepoName)
	return "", errors.New(errMessage)
}

func GetAliYunConnection() error {
	aliyunObj := new(AliyunApi)
	err := aliyunObj.createClient(
		tea.String(conf.Global.Aliyun.Access_key_id),
		tea.String(conf.Global.Aliyun.Access_key_secret),
		fmt.Sprintf("cr.cn-%s.aliyuncs.com", conf.Global.Aliyun.Region))
	if err != nil {
		return err
	}
	return nil
}

func CallAliYunApi(ctx context.Context, image map[string]string) string {
	defer func() {
		log.Logger.Debugf("[Github] 清理创建的仓库")
		err := githubApiSelf.DeleteRepo(ctx, conf.Global.Github.Namespace, image["image_name"])
		if err != nil {
			log.Logger.Debugf("[Github] 资源 %s:%s 清理失败：%s", image["image_name"], image["image_version"], err.Error())
		}
	}()

	log.Logger.Infof("[阿里云] 镜像 %s:%s 拉取进行中", image["image_name"], image["image_version"])
	aliyunObj := new(AliyunApi)
	aliyunObj.RepoName = image["image_name"]
	aliyunObj.RepoVersion = image["image_version"]

	err := aliyunObj.createRepository(ctx, image)
	if err != nil {
		log.Logger.Warnf("[阿里云] 镜像仓库 %s:%s 创建失败：%s", image["image_name"], image["image_version"], err.Error())
		return ""
	}

	imageExists, res, err := aliyunObj.createBuildRule()
	if err != nil {
		log.Logger.Warnf("[阿里云] 镜像仓库 %s:%s 构建规则创建失败：%s", image["image_name"], image["image_version"], err.Error())
		return ""
	}
	if imageExists {
		return res
	}

	imageURL, err := aliyunObj.startBuild(ctx, res, image)
	if err != nil {
		log.Logger.Warnf("[阿里云] 镜像 %s:%s 拉取失败：%s", image["image_name"], image["image_version"], err.Error())
		return ""
	}
	return imageURL
}
