package conf

import "encoding/json"

type RuleData struct {
	BuildRule BuildRule `json:"BuildRule"`
}

type BuildRule struct {
	DockerfileLocation string `json:"DockerfileLocation"`
	DockerfileName     string `json:"DockerfileName"`
	PushType           string `json:"PushType"`
	PushName           string `json:"PushName"`
	ImageTag           string `json:"ImageTag"`
	Tag                string `json:"Tag"`
}

// 创建阿里云仓库构建规则 request 参数
func GenRuleData(data RuleData) string {
	res, _ := json.Marshal(data)
	return string(res)
}
