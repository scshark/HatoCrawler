/**
 * @Author: scshark
 * @Description:
 * @File:  xgb
 * @Date: 1/9/23 1:31 PM
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
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
)

type Xgb struct{}

const (
	xgbGetInit = iota + 1
	xgbGetNew
	xgbGetIntervals
	xgbLiveUrlInit = "https://xuangubao.cn/live"
	xgbLiveUrlNew  = "https://api.xuangubao.cn/api/messages/live?limit=30"
	xgbLiveUrlMore = "https://baoer-api.xuangubao.cn/api/v6/message/newsflash?limit=20&cursor=%d&subj_ids=9,10,723,35,469,821&platform=pcweb"
)

func (crawler Xgb) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab.Xgb",
		Description: "选股宝",
	}
}

func (crawler Xgb) Get() error {

	logrus.Infof("%s 开始初始化",crawler.Config().Description)
	err := crawler.initXgbLives()
	if err != nil {
		logrus.Fatalf(" 选股宝初始化错误 err : %s", err)
	}
	logrus.Infof("%s 初始化成功",crawler.Config().Description)

	err = crawler.xgbCronCrawler()
	if err != nil {
		logrus.Fatalf(" 选股宝定时器启动失败 err : %s", err)
	}
	logrus.Infof("%s 定时采集启动 @every 1m runXgbLatest @every 3m runXgbIntervals  ",crawler.Config().Description)

	return err
}

func runXgbLatest() {
	_ = Xgb{}.getLiveData(xgbGetNew)
}
func runXgbIntervals() {
	_ = Xgb{}.getLiveData(xgbGetIntervals)
}

func (crawler Xgb) xgbCronCrawler() error {
	_cron := cron.New()
	err := _cron.AddFunc("@every 1m", runXgbLatest)
	err = _cron.AddFunc("@every 3m", runXgbIntervals)
	if err != nil {
		return err
	}
	_cron.Start()
	return nil
}
// 初始化 intervals
func (crawler Xgb) initXgbLives() error {
	// init 初始化
	logrus.Infof("%s 初始化获取数据",crawler.Config().Description)

	resp, err := GetUrlData(xgbLiveUrlInit,"html")
	if err != nil {
		logrus.Fatalf("选股宝初始化信息获取失败 GetData URL %s, error: %s", xgbLiveUrlInit,err)
	}
	logrus.Infof("%s 解析初始化数据 ",crawler.Config().Description)
	cursor, err := crawler.initRespParse(resp)

	if err != nil || cursor == 0 {
		logrus.Fatalf("选股宝初始化信息解析错误 GetData resp %s, error: %s", resp,err)
	}

	logrus.Infof("%s 初始化游标 ：%d",crawler.Config().Description,cursor)

	err = service.InitLivesIntervals(cursor, config.XgbLivesCrawler,0)
	if err != nil {
		logrus.Fatalf("选股宝游标初始化失败 error: %s",err)
	}
	return err
}

func (crawler Xgb) getLiveData(getType int) error {

	var url string
	switch getType {
	case xgbGetNew:
		logrus.Infof("%s 获取新数据 ：",crawler.Config().Description)
		url = xgbLiveUrlNew
	case xgbGetIntervals:
		// 获取cursor
		logrus.Infof("%s 获取区间数据 ：",crawler.Config().Description)

		cursor := service.GetLivesCursor(config.XgbLivesCrawler,0)
		if cursor == 0 {
			return errors.New("xgb游标is 0")
		}
		url = fmt.Sprintf(xgbLiveUrlMore, cursor)
	}

	resp, err := GetUrlData(url,"html")

	if err != nil {
		logrus.Errorf("XGB获取信息失败 url %s ,error %s",url,err)
	}

	var livesItem = make([]service.XgbLivesItems, 0)
	var next int64

	switch getType {
	case xgbGetNew:
		logrus.Infof("%s 解析新数据 ",crawler.Config().Description)
		livesItem, err = crawler.newRespParse(resp)
	case xgbGetIntervals:
		logrus.Infof("%s 解析区间数据 ",crawler.Config().Description)
		livesItem,next, err = crawler.intervalsRespParse(resp)
		if next <= 0 {
			logrus.Errorf("XGB获取游标信息失败 resp %s ,error %s",resp,err)
		}

	}

	if err != nil || len(livesItem) == 0 {
		logrus.Errorf("XGB获取信息失败 resp %s ,error %s",resp,err)
		return err
	}
	//
	logrus.Infof("%s 获取到 %d 条数据 ，开始保存数据",crawler.Config().Description,len(livesItem))

	// 数据保存
	err = service.CreateXgbLivesData(livesItem)
	if err != nil {
		logrus.Errorf("XGB信息保存失败 error %s",err)
		return err
	}

	if getType == xgbGetIntervals {
		// 处理 cursor
		logrus.Infof("%s 更新游标 ：%d",crawler.Config().Description,next)
		err = service.UpdateLivesCursor(config.XgbLivesCrawler, next,0)
		if err != nil {
			logrus.Errorf("XGB游标信息更新失败 error %s",err)
			return err
		}
	}
	return nil
}


func (crawler Xgb) initRespParse(resp string) (int64, error) {
	var err error
	var bodyString string
	// 去除STYLE
	re, _ := regexp.Compile(`\<style[\S\s]+?\</style\>`)
	bodyString = re.ReplaceAllString(resp, "")

	// //去除前面多余的信息
	re, _ = regexp.Compile(`\<!doctype html[\S\s]+?\<div infinite-scroll-disabled="`)
	bodyString = re.ReplaceAllString(bodyString, "")

	// 去除后面多余信息
	re, _ = regexp.Compile(`selectedFixedSubjectIds="[\S\s]+?\</html\>`)
	bodyString = re.ReplaceAllString(bodyString, "")


	// //去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	bodyString = re.ReplaceAllString(bodyString, "")

	re = regexp.MustCompile(`cursor="(\d*?)"`)
	result := re.FindStringSubmatch(strings.TrimSpace(bodyString))

	if len(result) < 1 {
		err = errors.New("XGB 初始化解析 respParse FindAllStringSubmatch error ")
		return 0, err
	}
	cursor, err := strconv.ParseInt(result[1], 10, 64)

	return cursor, err

}

//
func (crawler Xgb) newRespParse(resp string) ([]service.XgbLivesItems, error) {

	var err error
	resJson := gjson.Parse(resp)

	// 解析live
	respData := resJson.Get("Messages")

	if !respData.Exists() {
		err = errors.New("XGB 获取新信息 解析数据为空")
	}
	// slice of live list
	lives := make([]service.XgbLivesItems, 0)

	respData.ForEach(func(key, value gjson.Result) bool {

		var items service.XgbLivesItems

		value.ForEach(func(k, v gjson.Result) bool {

			switch k.Str {
			case "Id":
				items.Id = v.Int()
			case "Title":
				items.Title = v.Str
			case "Summary":
				items.Summary = v.Str
			case "Image":
				items.Image = v.Str
			case "OriginalUrl":
				items.OriginaUrl = v.Str
			case "CreatedAt":
				items.LiveCreatedAt = v.Int()
			case "Source":
				items.Source = v.Str
			case "LiveSubjects":
				// 单独处理
				sub := v.String()
				if sub == "" {
					break
				}
				subTitle := value.Get("LiveSubjects.#.Title")
				items.Tags = subTitle.String()
			case "Route":
				items.Uri = v.Str
			}
			return true
		})

		lives = append(lives, items)
		return true
	})
	return lives, err
}

func (crawler Xgb) intervalsRespParse(resp string) ([]service.XgbLivesItems, int64, error) {

	var err error
	resJson := gjson.Parse(resp)

	next := resJson.Get("data.next_cursor").Int()
	// 解析live
	respData := resJson.Get("data.messages")

	if !respData.Exists() {
		err = errors.New("XGB 获取区间数据 解析数据为空")
	}
	// slice of live list
	lives := make([]service.XgbLivesItems, 0)

	respData.ForEach(func(key, value gjson.Result) bool {

		var items service.XgbLivesItems

		value.ForEach(func(k, v gjson.Result) bool {

			switch k.Str {
			case "id":
				items.Id = v.Int()
			case "title":
				items.Title = v.Str
			case "summary":
				items.Summary = v.Str
			case "image":
				items.Image = v.Str
			case "created_at":
				items.LiveCreatedAt = v.Int()
			case "subj_ids":
				// 单独处理
				items.SubjIds = v.String()
			case "bkj_infos":
				// 单独处理
				bkj := v.String()
				if bkj == "" {
					break
				}
				bkjName := value.Get("bkj_infos.#.name")
				items.Tags = bkjName.String()
			case "route":
				items.Uri = v.Str
			}
			return true
		})

		lives = append(lives, items)
		return true
	})
	return lives,next, err
}
