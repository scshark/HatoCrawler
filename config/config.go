package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const Banner = `
   _____                    _           
 / ____|                  | |          
| |     _ __ __ ___      _| | ___ _ __ 
| |    | '__/ _  \ \ /\ / / |/ _ \ '__|
| |____| | | (_| |\ V  V /| |  __/ |   
 \_____|_|  \__,_| \_/\_/ |_|\___|_|   

`

// 全局Config
var Cfg *Config
var (
	Test       bool
	Version    bool
	Help       bool
	Generate   bool
	ConfigFile string

	GITHUB    string = "https://github.com/scshark/HatoCrawler"
	TAG       string = "v0.1"
	GOVERSION string = "go1.19.3"

)



func DefaultConfig() Config {
	return Config{
		ChromeDriver: "./chromedriver/linux64",
		Proxy: ProxyStruct{
			ProxyUrl:            "http://127.0.0.1:7890",
			CrawlerProxyEnabled: false,
			BotProxyEnabled:     false,
		},
		Cron: CronStruct{
			Enabled: false,
			Time:    11,
		},
		Api: ApiStruct{
			Enabled: false,
			Debug:   false,
			Host:    "127.0.0.1",
			Port:    8080,
			Auth:    "auth_key_here",
		},
		Crawler: CrawlerStruct{
			EdgeForum:   EdgeForumStruct{Enabled: false},
			XianZhi:     XianZhiStruct{Enabled: false, UseChromeDriver: true, CustomRSSURL: ""},
			SeebugPaper: SeebugPaperStruct{Enabled: false},
			Anquanke:    AnquankeStruct{Enabled: false},
			Tttang:      TttangStruct{Enabled: false},
			QiAnXin:     QiAnXinStruct{Enabled: false},
			// DongJian:    DongJianStruct{Enabled: false},
			Lab: LabStruct{
				Enabled:     true,
				NoahLab:     NoahLabStruct{Enabled: false},
				Blog360:     Blog360Struct{Enabled: false},
				Nsfocus:     NsfocusStruct{Enabled: false},
				Xlab:        XlabStruct{Enabled: false},
				AlphaLab:    AlphaLabStruct{Enabled: false},
				Netlab:      NetlabStruct{Enabled: false},
				RiskivyBlog: RiskivyBlogStruct{Enabled: false},
				TSRCBlog:    TSRCBlogStruct{Enabled: false},
				X1cT34m:     X1cT34mStruct{Enabled: false},
				Jinse:       JinseStruct{Enabled: true},
				WallStreet:  WallStreetStruct{Enabled: true},
				Twitter:     TwitterStruct{Enabled: true, ScreenName: []string{}},
			},
		},
		Bot: BotStruct{
			WecomBot: WecomBotStruct{
				Enabled: false,
				Key:     "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				Timeout: 2,
			},
			FeishuBot: FeishuBotStruct{
				Enabled: false,
				Key:     "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				Timeout: 2,
			},
			DingBot: DingBotStruct{
				Enabled: false,
				Token:   "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Timeout: 2,
			},
			HexQBot: HexQBotStruct{
				Enabled: false,
				Api:     "http://xxxxxx.com/send",
				QQGroup: 000000000,
				Key:     "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				Timeout: 2,
			},
			ServerChan: ServerChanStruct{
				Enabled: false,
				SendKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Timeout: 2,
			},
			WgpSecBot: WgpSecBotStruct{
				Enabled: false,
				Key:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Timeout: 2,
			},
		},
		Database: DatabaseStruct{
			TablePrefix:  "ht_",
			LogLevel:     "Error",
			Host:         "",
			UserName:     "",
			Password:     "",
			DBName:       "hato_crawler",
			Charset:      "utf8mb4",
			ParseTime:    true,
			MaxIdleConns: 10,
			MaxOpenConns: 30,
		},
		Redis: RedisStruct{
			Host:     "127.0.0.1",
			Password: "",
			DB:       0,
		},
		Logger: LoggerStruct{
			Level: "debug",
		},
		LoggerZinc: LoggerZincStruct{
			Host:     "",
			Index:    "",
			User:     "",
			Password: "",
			Secure:   false,
		},
	}
}
func (s *LoggerZincStruct) Endpoint() string {
	return endpoint(s.Host, s.Secure)
}

func configToYaml() string {
	b, err := yaml.Marshal(DefaultConfig())
	if err != nil {
		log.Fatalf("unable to marshal config to yaml: %s", err.Error())
	}
	return string(b)
}

func ConfigInit() {
	log.SetPrefix("[!] ")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(ConfigFile)
	// 判断config文件是否存在
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		if Generate {
			f, err := os.Create(ConfigFile)
			if err != nil {
				log.Fatalf("create config file error: %s\n", err.Error())
			}
			defer f.Close()

			_, err = f.WriteString(configToYaml())
			if err != nil {
				log.Fatalf("write config file error: %s\n", err.Error())
			}
			f.Sync()
			fmt.Println("[*] The configuration file has been initialized.")
			os.Exit(0)
		} else {
			fmt.Println("[!] The configuration file does not exist, please use `-init`")
			os.Exit(0)
		}
	} else {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("read config file error: %s\n", err.Error())
		}

		err = viper.Unmarshal(&Cfg)
		if err != nil {
			log.Fatalf("unmarshal config error: %s\n", err.Error())
		}
		fmt.Printf("[*] load config success!\n\n")
	}
}

func (s *LoggerStruct) logLevel() logrus.Level {
	switch strings.ToLower(Cfg.Logger.Level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.ErrorLevel
	}
}

func endpoint(host string, secure bool) string {
	schema := "http"
	if secure {
		schema = "https"
	}
	return schema + "://" + host
}
