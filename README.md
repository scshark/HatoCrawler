<h1 align="center">
HatoCrawler
</h1>






## Table of Contents

- [Usage](#usage) 
- [守护进程](#守护进程配置)
- [API](#api) 
- [ChromeDriver](#chromedriver) 
- [微信或QQ推送群](#微信或QQ推送群)
- [Install](#install)
- [Config](#config)
- [License](#license)


### Usage

程序使用yml格式的配置文件，第一次使用时请使用`-init`参数在当前文件夹生成默认配置文件，在配置文件中设置爬取的网站和推送机器人相关配置，目前包括在内的网站和推送的机器人在[Features](#features)中可以查看，可以设置每日推送的整点时间以及是否开启API。

```text

  _____                    _           
 / ____|                  | |          
| |     _ __ __ ___      _| | ___ _ __ 
| |    | '__/ _  \ \ /\ / / |/ _ \ '__|
| |____| | | (_| |\ V  V /| |  __/ |   
 \_____|_|  \__,_| \_/\_/ |_|\___|_|   							  
HatoCrawler dev

Options:
  -c file
    	the config file to be used, or generate a config file with the specified name with -init (default "config.yml")
  -help
    	print help info
  -init
    	generate a config file
  -test
    	stop after running once
  -version
    	print version info

```

- 使用`-h/-help`查看详细命令
- 使用`-c`指定使用的配置文件，或者在生成配置文件时配合`-init`生成指定文件名的配置文件
- 使用`-test`参数执行一次程序后退出
- 使用`-version`输出详细版本信息

如果开启了定时任务（Cron），程序使用定时任务每天根据设置好的时间整点自动运行，编辑好相关配置后后台运行即可。

简单运行命令：

```sh
$ nohup ./HatoCrawler >> run.log 2>&1 &
```

或者使用screen

```sh
$ screen ./HatoCrawler
$ ctrl a+d / control a+d # 回到主会话
```

如果长期使用，建议配置[守护进程](#守护进程配置)。
### 守护进程配置

首先执行`vim /etc/systemd/system/HatoCrawler.service`输入以下内容：

```
[Unit]
Description=HatoCrawler
After=network.target
 
[Service]
Type=simple
WorkingDirectory=<HatoCrawler Path>
ExecStart=<HatoCrawler Path>/HatoCrawler -c config.yml
Restart=on-failure
 
[Install]
WantedBy=multi-user.target
```

其中`<HatoCrawler Path>`为HatoCrawler可执行文件存放的路径。

保存后执行`systemctl daemon-reload`，现在你就可以使用以下命令来管理程序了：

- 启动: systemctl start HatoCrawler
- 关闭: systemctl stop HatoCrawler
- 自启: systemctl enable HatoCrawler
- 状态: systemctl status HatoCrawler
- 重启: systemctl restart HatoCrawler
- **查看日志**: journalctl -u HatoCrawler



程序旨在帮助安全研究者自动化获取每日更新的安全文章，适用于每日安全日报推送，爬取的安全社区网站范围和支持推送的机器人持续增加中，欢迎在[issues](https://github.com/Le0nHato/HatoCrawler/issues)中提供宝贵的建议。


:rocket: 目前 HatoCrawler 已在MacOS Apple silicon 、Ubuntu 20.04运行测试通过。
### API

HatoCrawler提供了Web API，配合其他工具可以主动调用API进行爬取或推送。

- [API文档](https://www.apifox.cn/apidoc/shared-b613c4fc-56a6-4724-831f-4c1ac5547ab5)
- 注意请求API需要带上Authorization头，在配置文件中配置`auth`值

### ChromeDriver

ChromeDriver镜像站：http://npm.taobao.org/mirrors/chromedriver/


- Windows和Mac用户在[下载Chrome](https://www.google.cn/chrome/)并安装后，下载对应chrome版本的ChromeDriver并在配置文件`config.yml`中指定ChromeDriver的路径
- Linux用户在下载Chrome（链接如下）并安装后，同上编辑配置文件
    - [Debian/Ubuntu(64位.deb)](https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb)
    - [Fedora/openSUSE(64位.rpm)](https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm)


> Chrome浏览器可以访问`chrome://version/`查看版本

> 命令行可以使用`google-chrome-stable --version`查看版本

### 微信或QQ推送群

如果不想自己配置环境，只想获取每日推送，可以扫码加推送群：

如果微信群二维码失效或者人数已满，可以添加微信号：WgpHatoBot，然后私聊发送 HatoCrawler 进群。




支持的推送机器人列表：

- [x] [企业微信群机器人](https://work.weixin.qq.com/api/doc/90000/90136/91770)
- [x] [HexQBot](https://github.com/Am473ur/HexQBot) (QQ群机器人 自建)
- [x] [Server酱](https://sct.ftqq.com/)
- [x] [飞书群机器人](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- [x] [钉钉群机器人](https://open.dingtalk.com/document/robots/custom-robot-access)
- [x] [WgpHatoBot](https://bot.wgpHato.org)
- [ ] [pushplus](http://pushplus.hxtrip.com/)

## Install

你可以在[Releases](https://github.com/Le0nHato/HatoCrawler/releases)下载最新的HatoCrawler。


或者从源码编译：


```sh
$ git clone https://github.com/Le0nHato/HatoCrawler.git
$ cd HatoCrawler
$ go build .
```


## Config
`config.yml`配置文件模板注释：

```yml
# 设置Selenium使用的ChromeDriver路径，支持相对路径或绝对路径（如果不爬取先知社区可以不用设置）
ChromeDriver: ./chromedriver/linux64

Proxy:
  ProxyUrl: http://127.0.0.1:7890 # 代理地址，支持http/https/socks协议
  CrawlerProxyEnabled: false # 是否开启爬虫代理
  BotProxyEnabled: false # 是否开启请求机器人代理

Cron:
  enabled: false # 是否开启定时任务，开启后每天按照指定的时间爬取并推送
  time: 11 # 设置定时任务每天整点爬取推送时间，范围 0 ~ 23（整数）

Api:
  enabled: false # 是否开启API
  debug: false # 是否开启Gin-DEBUG模式
  host: 127.0.0.1
  port: 8080
  auth: auth_key_here # 请求api需要带上Authorization头

Crawler:
  # 棱角社区
  # https://forum.ywhack.com/forum-59-1.html
  EdgeForum:
    enabled: false
  # 先知安全技术社区
  # https://xz.aliyun.com/
  XianZhi:
    enabled: false
    UHatohromeDriver: true # 是否使用selenium调用浏览器爬取，设置为true需要指定ChromeDriver地址，为false需要指定没有反爬措施的自定义网址CustomRSSURL
    CustomRSSURL: ""
  # SeebugPaper（知道创宇404实验室）
  # https://paper.seebug.org/
  SeebugPaper:
    enabled: false
  # 安全客
  # https://www.anquanke.com/
  Anquanke:
    enabled: false
  # 跳跳糖
  # http://tttang.com/
  Tttang:
    enabled: false
  # 奇安信攻防社区
  # https://forum.butian.net/community/all/newest
  QiAnXin:
    enabled: false
  # 洞见微信聚合
  # http://wechat.doonHato.com/
  # DongJian:
  #   enabled: false
  Lab:
    enabled: true # 是否开启各大实验室文章爬取
    NoahLab:
      enabled: true
    Blog360:
      enabled: true
    Nsfocus:
      enabled: true
    Xlab:
      enabled: true
    AlphaLab:
      enabled: true
    Netlab:
      enabled: true
    RiskivyBlog:
      enabled: true
    TSRCBlog:
      enabled: true
    X1cT34m:
      enabled: true
Bot:
  # 企业微信群机器人
  # https://work.weixin.qq.com/api/doc/90000/90136/91770
  WecomBot:
    enabled: false
    key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    timeout: 2
  # 飞书群机器人
  # https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN
  FeishuBot:
    enabled: false
    key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    timeout: 2
  # 钉钉群机器人
  # https://open.dingtalk.com/document/robots/custom-robot-access
  DingBot:
    enabled: false
    token: xxxxxxxxxxxxxxxxxxxx
    timeout: 2
  # HexQBot
  # https://github.com/Am473ur/HexQBot
  HexQBot:
    enabled: false
    api: http://xxxxxx.com/send
    qqgroup: 0
    key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    timeout: 2
  # Server酱
  # https://sct.ftqq.com/
  ServerChan:
    enabled: false
    sendkey: xxxxxxxxxxxxxxxxxxxx
    timeout: 2
  # WgpHatoBot
  # https://bot.wgpHato.org/
  WgpHatoBot:
    enabled: false
    key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    timeout: 2

```



   
## License


[Apache License 2.0](https://github.com/scshark/HatoCrawler/blob/main/LICENSE)
