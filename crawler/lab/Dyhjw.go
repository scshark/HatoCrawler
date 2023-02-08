/**
 * @Author: scshark
 * @Description:
 * @File:  Dyhjw
 * @Date: 1/5/23 3:48 PM
 */
package lab

import (
	"SecCrawler/config"
	"SecCrawler/internal/service"
	"SecCrawler/register"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Dyhjw struct{}

const (
	hjGetInit = iota + 1
	hjGetNew
	hjGetIntervals
	hjLiveUrlNew  = "http://www.dyhjw.com/kuaixun/"
	hjLiveUrlMore = "http://www.dyhjw.com/kuaixun/more?lastid=%s"
)

func (crawler Dyhjw) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab.Dyhjw",
		Description: "第一黄金网",
	}
}

func (crawler Dyhjw) Get() error {

	// 初始化信息
	logrus.Infof("%s 开始初始化",crawler.Config().Description)
	err := crawler.getLiveData(jinseInit)
	if err != nil {
		logrus.Fatalf(" 第一黄金网获取数据初始化失败 err : %s", err)
	}
	logrus.Infof("%s 初始化成功",crawler.Config().Description)

	err = crawler.dhCronCrawler()
	if err != nil {
		logrus.Fatalf("第一黄金网定时器启动失败 %s", err)
	}
	logrus.Infof("%s 定时采集启动 @every 1m runDhLates  @every 5m runDhIntervals",crawler.Config().Description)

	return err
}
func runDhLatest() {
	_ = Dyhjw{}.getLiveData(hjGetNew)
}
func runDhIntervals() {
	_ = Dyhjw{}.getLiveData(hjGetIntervals)
}
func (crawler Dyhjw) dhCronCrawler() error {

	_cron := cron.New()
	err := _cron.AddFunc("@every 1m", runDhLatest)
	err = _cron.AddFunc("@every 5m", runDhIntervals)
	_cron.Start()
	return err
}

func (crawler Dyhjw) getLiveData(getType int) error {

	var url string
	switch getType {
	case hjGetInit,hjGetNew:
		logrus.Infof("%s 开始新数据获取",crawler.Config().Description)
		url = hjLiveUrlNew
	case hjGetIntervals:
		// 获取cursor
		logrus.Infof("%s 开始区间数据获取",crawler.Config().Description)
		cursor := service.GetLivesCursor(config.DyhjwLivesCrawler, 0)
		if cursor == 0 {
			return errors.New("第一黄金网获取游标 为 0，无法进行数据获取")
		}
		url = fmt.Sprintf(hjLiveUrlMore, strconv.FormatInt(cursor, 10))
	}

	// 获取数据
	resp, err := GetUrlData(url, "html")

	if err != nil {
		logrus.Errorf("第一黄金网数据获取失败 url %s, error %s", url, err)
		return err
	}
	logrus.Infof("%s 开始解析数据",crawler.Config().Description)

	livesItem, err := crawler.respParse(resp)

	logrus.Infof("%s 解析数据成功",crawler.Config().Description)


	if err != nil || len(livesItem) == 0 {
		logrus.Errorf("第一黄金网数据解析失败 resp %v, error %s", resp, err)
		return err
	}
	logrus.Infof("%s 获取到数据 %d条,开始保存数据",crawler.Config().Description,len(livesItem))

	// 保存数据
	err = service.CreateDyhjwLivesData(livesItem)
	if err != nil {
		logrus.Errorf("第一黄金网数据保存失败 error %s", err)
		return err
	}

	switch getType {
	case wsGetInit:
		// 获取id
		// 获取最后一个id
		logrus.Infof("%s 开始初始化游标 %d ",crawler.Config().Description,livesItem[len(livesItem)-1].Id)

		cursor, _ := strconv.ParseInt(livesItem[len(livesItem)-1].Id, 10, 64)
		err = service.InitLivesIntervals(cursor, config.DyhjwLivesCrawler, 0)
		if err != nil {
			logrus.Fatalf("第一黄金网初始化游标失败 %s", err)
		}
	case wsGetIntervals:
		logrus.Infof("%s 更新游标 ： %d",crawler.Config().Description,livesItem[len(livesItem)-1].Id)

		cursor, _ := strconv.ParseInt(livesItem[len(livesItem)-1].Id, 10, 64)
		err = service.UpdateLivesCursor(config.DyhjwLivesCrawler, cursor, 0)
		if err != nil {
			logrus.Fatalf("第一黄金网更新游标失败 %s", err)
		}

	}

	return err
}

func (crawler Dyhjw) respParse(resp string) ([]service.DyhjwLivesItems, error) {

	var err error
	var bodyString string
	// 去除STYLE
	re, _ := regexp.Compile(`\<style[\S\s]+?\</style\>`)
	bodyString = re.ReplaceAllString(resp, "")
	//
	// 去除SCRIPT
	// re, _ = regexp.Compile(`\<script[\S\s]+?\</script\>`)
	// bodyString = re.ReplaceAllString(bodyString, "")

	// //去除head
	// re, _ = regexp.Compile(`\<head[\S\s]+?\</head\>`)
	// bodyString = re.ReplaceAllString(bodyString, "")

	// //去除前面多余的信息
	re, _ = regexp.Compile(`\<!DOCTYPE HTML[\S\s]+?id="kxlist"\>`)
	bodyString = re.ReplaceAllString(bodyString, "")

	// //去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	bodyString = re.ReplaceAllString(bodyString, "")

	// 去除后面多余信息
	re, _ = regexp.Compile(`\<\/ul><div class="more_news[\S\s]+?\</html\>`)
	bodyString = re.ReplaceAllString(bodyString, "")

	// 继续优化去处多余信息
	re, _ = regexp.Compile(`\<table cellpadding="[\S\s]+?class="kx_title"\>`)
	bodyString = re.ReplaceAllString(bodyString, "")

	re, _ = regexp.Compile(`\<\/p\>\<[\S\s]+?table\>`)
	bodyString = re.ReplaceAllString(bodyString, "")

	re = regexp.MustCompile(`id="(\d*?)">(.*?)\<\/li\>`)
	result := re.FindAllStringSubmatch(strings.TrimSpace(bodyString), -1)

	if len(result) == 0 {
		err = errors.New("respParse FindAllStringSubmatch error")
		return nil, err
	}
	var respLives = make([]service.DyhjwLivesItems, 0)
	for _, items := range result {
		if len(items) < 3 {
			continue
		}
		d := service.DyhjwLivesItems{}
		for k, v := range items {

			switch k {
			case 1:
				timeZone := time.FixedZone("CST", 8*3600)
				t := fmt.Sprintf("%s %s", v[:8], v[8:14])
				displayTime, err := time.ParseInLocation("20060102 150405", t, timeZone)
				if err != nil {
					logrus.Errorf("dh 解析 items time Parse error %s", err)
					continue
				}

				d.DisplayTime = displayTime.Unix()
				d.Id = v
			case 2:
				d.Content = v
			}
		}
		respLives = append(respLives, d)
	}

	return respLives, err
}
