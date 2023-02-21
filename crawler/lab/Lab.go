package lab

import (
	. "HatoCrawler/config"
	"HatoCrawler/register"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Lab struct{}

func (crawler Lab) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab",
		Description: "实验室文章",
	}
}

// Get 获取 Lab 前24小时内文章。
func (crawler Lab) Get() error {
	var resultSlice [][]string

	// if Cfg.Crawler.Lab.NoahLab.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, NoahLab{})
	// }
	// if Cfg.Crawler.Lab.Blog360.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, Blog360{})
	// }
	// if Cfg.Crawler.Lab.Nsfocus.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, Nsfocus{})
	// }
	// if Cfg.Crawler.Lab.Xlab.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, Xlab{})
	// }
	// if Cfg.Crawler.Lab.AlphaLab.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, AlphaLab{})
	// }
	// if Cfg.Crawler.Lab.Netlab.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, Netlab{})
	// }
	// if Cfg.Crawler.Lab.RiskivyBlog.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, RiskivyBlog{})
	// }
	// if Cfg.Crawler.Lab.TSRCBlog.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, TSRCBlog{})
	// }
	// if Cfg.Crawler.Lab.X1cT34m.Enabled {
	// 	resultSlice = tmpCrawler(resultSlice, X1cT34m{})
	// }


	if Cfg.Crawler.Lab.Jinse.Enabled {
		resultSlice = tmpCrawler(resultSlice, Jinse{})
	}
	if Cfg.Crawler.Lab.WallStreet.Enabled {
		resultSlice = tmpCrawler(resultSlice, WallStreet{})
	}
	if Cfg.Crawler.Lab.Dyhjw.Enabled {
		resultSlice = tmpCrawler(resultSlice,Dyhjw{})
	}
	if Cfg.Crawler.Lab.Xgb.Enabled {
		resultSlice = tmpCrawler(resultSlice,Xgb{})
	}
	if Cfg.Crawler.Lab.Twitter.Enabled {
		resultSlice = tmpCrawler(resultSlice,Twitter{ ScreenName:Cfg.Crawler.Lab.Twitter.ScreenName})
	}
	if len(resultSlice) == 0 {
		return  errors.New("no records in the last 24 hours")
	}
	return  nil
}

func tmpCrawler(s [][]string, crawler register.Crawler) [][]string {
	err := crawler.Get()
	if err != nil {
		logrus.Errorf("crawl [%s] error: %s\n\n", crawler.Config().Name, err.Error())
	}
	logrus.Infof("%s采集器： 启动成功",crawler.Config().Description)
	logrus.Info("采集器启动 time sleep 10 second")
	time.Sleep(10 * time.Second)
	return s
}
