# 1.imagectl 是什么？

使用 kubernetes 时，或许经常出现这样的场景：“Pod 创建过程失败，并处于 ErrImagePull 状态”，这是由于部分镜像网站被 GFW 阻挡，导致无法拉取镜像，如常见的 gcr.io 等等

这时，就需要将镜像同步至阿里云镜像 ACR 以供下载
imagectl 则是将此**过程自动化**的一个小工具

支持多协程拉取镜像

# 2.版本说明

基于 Golang 1.17.2 开发

# 3.安装

最新二进制程序文件可访问 [releases](https://github.com/hongqchen/imagectl/releases) 获取

# 4.使用说明

## 4.1 Github

- 登录 [github](https://github.com/)
- 点击右上角的头像，并选中 Settings 选项
- 左侧菜单栏选中 Developer settings
- Personal access tokens 选中 Generate new token

## 4.2 阿里云

- [登录](https://usercenter.console.aliyun.com/) 新建 access token 并记录
- [登录阿里云镜像服务](https://cr.console.aliyun.com/cn-chengdu/instance/namespaces)创建命名空间，请自行切换地域

## 4.3 imagectl

### 配置文件填写

> 将刚才记录的信息，填写至 config/base.conf 文件中

**配置文件字段说明**

```
[github]
access_token = ""   // 填写 github access token
namespace = ""      // 填写 github 名称空间

[aliyun]
access_key_id = ""  // 填写阿里云提供的 access key ID
access_key_secret = "" // 填写阿里云提供的 access key secret
namespace = "" // 填写在阿里云镜像服务网站里创建的镜像名称空间
repo_type = ""    // 镜像访问方式，PUBLIC 可直接 docker pull；PRIVATE 需要 docker login 后，才能 pull 镜像
region = "" // 填写上述命名空间所在的地域名称，如 成都:chengdu,南京:nanjing,上海:shanghai 等
```

### 程序运行

可同时拉取多个镜像，中间用空格分割，--debug 选项查看详细拉取步骤
只测试了 intel 版本 mac，linux 和 windows 未测试

```
linux,macos:
imagectl sync 'image_url1'
imagectl sync 'image_url1' 'image_url2'
imagectl sync 'image_url1' 'image_url2' --debug
imagectl --help 可查看帮助信息

windows:
imagectl.exe sync 'image_url1'
其它选项同上
```

**构建速度快与慢取决于阿里云的速度，与程序无关，如果镜像构建失败，大概率镜像不存在或阿里云镜像站本身存在问题，可尝试重新构建**

### 运行效果

**镜像同步成功**
```
ChenHQ-MacBook-Pro:~ chenhongquan$ imagectl sync  'k8s.gcr.io/prometheus-adapter/prometheus-adapter:v0.9.1' 'k8s.gcr.io/metrics-server/metrics-server:v0.6.1'
2022-04-09 23:33:37.766	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 metrics-server:v0.6.1 拉取进行中
2022-04-09 23:33:37.766	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 prometheus-adapter:v0.9.1 拉取进行中
2022-04-09 23:34:20.924	[INFO]	cmd/start.go:55	[阿里云] 镜像地址:
registry.cn-chengdu.aliyuncs.com/hongqchen/prometheus-adapter:v0.9.1
registry.cn-chengdu.aliyuncs.com/hongqchen/metrics-server:v0.6.1
```

**部分镜像同步成功，部分失败**
```
ChenHQ-MacBook-Pro:~ chenhongquan$ ./imagectl sync  'k8s.gcr.io/prometheus-adapter/prometheus-adapter:v0.9.1' 'k8s.gcr.io/metrics-server/metrics-server:v0.6.1111'
2022-04-09 23:34:55.193	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 metrics-server:v0.6.1111 拉取进行中
2022-04-09 23:34:55.193	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 prometheus-adapter:v0.9.1 拉取进行中
2022-04-09 23:36:06.005	[WARN]	aliyunApiSelf/aliyun.go:318	[阿里云] 镜像 metrics-server:v0.6.1111 拉取失败：错误原因登录 https://cr.console.aliyun.com/repository/chengdu/hongqchen/metrics-server/build 查看
2022-04-09 23:36:06.903	[INFO]	cmd/start.go:55	[阿里云] 镜像地址:
registry.cn-chengdu.aliyuncs.com/hongqchen/prometheus-adapter:v0.9.1
```

**开启 debug 模式**
```
ChenHQ-MacBook-Pro:imagectl-macos-amd64 chenhongquan$ ./imagectl sync  'k8s.gcr.io/prometheus-adapter/prometheus-adapter:v0.9.1' 'k8s.gcr.io/metrics-server/metrics-server:v0.6.1111'  --debug
2022-04-10 00:43:53.483	[DEBUG]	cmd/start.go:16	参数校验中...
2022-04-10 00:43:53.483	[DEBUG]	cmd/start.go:22	参数输入正确.
2022-04-10 00:43:53.483	[DEBUG]	cmd/start.go:24	配置文件加载中...
2022-04-10 00:43:53.484	[DEBUG]	cmd/start.go:30	配置文件加载完成.
2022-04-10 00:43:53.484	[DEBUG]	cmd/start.go:32	[Github] 建立连接中...
2022-04-10 00:43:53.484	[DEBUG]	cmd/start.go:35	[阿里云] 建立连接中...
2022-04-10 00:43:53.484	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 metrics-server:v0.6.1111 拉取进行中
2022-04-10 00:43:53.484	[INFO]	aliyunApiSelf/aliyun.go:296	[阿里云] 镜像 prometheus-adapter:v0.9.1 拉取进行中
2022-04-10 00:43:53.484	[DEBUG]	aliyunApiSelf/aliyun.go:77	[阿里云] 判断镜像仓库 metrics-server:v0.6.1111 是否存在
2022-04-10 00:43:53.484	[DEBUG]	aliyunApiSelf/aliyun.go:77	[阿里云] 判断镜像仓库 prometheus-adapter:v0.9.1 是否存在
2022-04-10 00:43:54.074	[DEBUG]	aliyunApiSelf/aliyun.go:86	[阿里云] 镜像仓库 prometheus-adapter:v0.9.1 创建中
2022-04-10 00:43:54.074	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 prometheus-adapter 是否存在
2022-04-10 00:43:54.096	[DEBUG]	aliyunApiSelf/aliyun.go:86	[阿里云] 镜像仓库 metrics-server:v0.6.1111 创建中
2022-04-10 00:43:54.096	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 metrics-server 是否存在
2022-04-10 00:43:54.641	[DEBUG]	githubApiSelf/github.go:100	[Github] 镜像仓库 metrics-server:v0.6.1111 创建中
2022-04-10 00:43:54.642	[DEBUG]	githubApiSelf/github.go:100	[Github] 镜像仓库 prometheus-adapter:v0.9.1 创建中
2022-04-10 00:43:56.772	[DEBUG]	aliyunApiSelf/aliyun.go:139	[阿里云] 判断镜像 metrics-server:v0.6.1111 是否存在
2022-04-10 00:43:56.779	[DEBUG]	aliyunApiSelf/aliyun.go:139	[阿里云] 判断镜像 prometheus-adapter:v0.9.1 是否存在
2022-04-10 00:43:56.967	[DEBUG]	aliyunApiSelf/aliyun.go:151	[阿里云] 查询镜像 prometheus-adapter:v0.9.1 构建规则是否存在
2022-04-10 00:43:56.967	[DEBUG]	aliyunApiSelf/aliyun.go:151	[阿里云] 查询镜像 metrics-server:v0.6.1111 构建规则是否存在
2022-04-10 00:43:57.135	[DEBUG]	aliyunApiSelf/aliyun.go:161	[阿里云] 镜像 prometheus-adapter:v0.9.1 构建规则不存在，规则创建中
2022-04-10 00:43:57.135	[DEBUG]	aliyunApiSelf/aliyun.go:161	[阿里云] 镜像 metrics-server:v0.6.1111 构建规则不存在，规则创建中
2022-04-10 00:43:57.304	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 prometheus-adapter 是否存在
2022-04-10 00:43:57.324	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 metrics-server 是否存在
2022-04-10 00:43:58.360	[DEBUG]	aliyunApiSelf/aliyun.go:199	[阿里云] 镜像 prometheus-adapter:v0.9.1 构建状态查询
2022-04-10 00:43:58.574	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: PENDING，等待 5s 后再次查询任务状态
2022-04-10 00:43:58.692	[DEBUG]	aliyunApiSelf/aliyun.go:199	[阿里云] 镜像 metrics-server:v0.6.1111 构建状态查询
2022-04-10 00:43:58.839	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] metrics-server:v0.6.1111 构建状态: PENDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:03.754	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: PENDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:03.994	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] metrics-server:v0.6.1111 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:50.190	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:50.459	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] metrics-server:v0.6.1111 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:55.349	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:44:55.669	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] metrics-server:v0.6.1111 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:45:00.499	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:45:00.814	[WARN]	aliyunApiSelf/aliyun.go:318	[阿里云] 镜像 metrics-server:v0.6.1111 拉取失败：错误原因登录 https://cr.console.aliyun.com/repository/chengdu/hongqchen/metrics-server/build 查看
2022-04-10 00:45:00.814	[DEBUG]	aliyunApiSelf/aliyun.go:289	[Github] 清理创建的仓库
2022-04-10 00:45:00.814	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 metrics-server 是否存在
2022-04-10 00:45:05.673	[DEBUG]	aliyunApiSelf/aliyun.go:262	[阿里云] prometheus-adapter:v0.9.1 构建状态: BUILDING，等待 5s 后再次查询任务状态
2022-04-10 00:45:10.831	[DEBUG]	aliyunApiSelf/aliyun.go:289	[Github] 清理创建的仓库
2022-04-10 00:45:10.831	[DEBUG]	githubApiSelf/github.go:65	[Github] 判断镜像仓库 prometheus-adapter 是否存在
2022-04-10 00:45:11.882	[INFO]	cmd/start.go:55	[阿里云] 镜像地址:
registry.cn-chengdu.aliyuncs.com/hongqchen/prometheus-adapter:v0.9.1
```

