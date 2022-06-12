package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

var Global *BaseConfig

type BaseConfig struct {
	Github *GithubConfig `toml:"github"`
	Aliyun *AliyunConfig `toml:"aliyun"`
}

type GithubConfig struct {
	Access_token string `toml:"access_token"`
	Namespace    string `toml:"namespace"`
}

type AliyunConfig struct {
	Access_key_id     string `toml:"access_key_id"`
	Access_key_secret string `toml:"access_key_secret"`
	Namespace         string `toml:"namespace"`
	Repo_type         string `toml:"repo_type"`
	Region            string `toml:"region"`
}

func LoadConfig() error {
	var baseData BaseConfig
	ex, _ := os.Executable()
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config/base.conf", filepath.Dir(ex)), &baseData); err != nil {
		return err
	}
	Global = &baseData
	return nil
}
