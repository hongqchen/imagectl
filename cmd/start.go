package cmd

import (
	"context"
	"fmt"
	"github.com/hongqchen/imagectl/aliyunApiSelf"
	"github.com/hongqchen/imagectl/githubApiSelf"
	"github.com/hongqchen/imagectl/utils/conf"
	"github.com/hongqchen/imagectl/utils/log"
	"sync"
)

var wg sync.WaitGroup

func Start(imageUrls []string) {
	log.Logger.Debug("参数校验中...")
	images, err := conf.ParseImageUrls(imageUrls)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	log.Logger.Debug("参数输入正确.")

	log.Logger.Debug("配置文件加载中...")
	err = conf.LoadConfig()
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	log.Logger.Debug("配置文件加载完成.")

	log.Logger.Debug("[Github] 建立连接中...")
	githubApiSelf.GetGithubConnection()

	log.Logger.Debug("[阿里云] 建立连接中...")
	err = aliyunApiSelf.GetAliYunConnection()
	if err != nil {
		log.Logger.Warn("[阿里云] 连接失败!")
		return
	}

	imageUrlsChan := make(chan string, len(images))

	for _, image := range images {
		wg.Add(1)
		go func(image map[string]string, wg *sync.WaitGroup) {
			defer wg.Done()
			imageUrl := aliyunApiSelf.CallAliYunApi(context.Background(), image)
			imageUrlsChan <- imageUrl
		}(image, &wg)

	}
	wg.Wait()

	log.Logger.Info("[阿里云] 镜像地址:")
	for i := 1; i <= len(images); i++ {
		fmt.Printf(<-imageUrlsChan)
	}
	defer close(imageUrlsChan)
}
