package conf

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func ParseImageUrls(urls []string) ([]map[string]string, error) {
	parsedUrls := []map[string]string{}
	for _, url := range urls {
		urlSplit := strings.Split(url, "/")

		tmpVersion := []string{}
		switch len(urlSplit) {
		case 2:
			tmpVersion = strings.Split(urlSplit[1], ":")
		case 3:
			tmpVersion = strings.Split(urlSplit[2], ":")
		default:
			return nil, errors.New(fmt.Sprintf("镜像地址: [%s] 格式错误!", url))
		}

		if len(tmpVersion) != 2 || tmpVersion[1] == "" {
			return nil, errors.New(fmt.Sprintf("镜像地址: [%s] 格式错误!", url))
		}

		parsedUrl := map[string]string{
			"image_name":    tmpVersion[0],
			"image_version": tmpVersion[1],
			"image_url":     url,
		}

		parsedUrls = append(parsedUrls, parsedUrl)
	}
	return parsedUrls, nil
}
