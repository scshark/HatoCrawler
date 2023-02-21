/**
 * @Author: scshark
 * @Description:
 * @File:  WallStreet
 * @Date: 12/29/22 2:00 PM
 */
package lab

import (
	"HatoCrawler/config"
	"HatoCrawler/internal/service"
	"HatoCrawler/register"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
)

type WallStreet struct{}

const (
	wsGetInit = iota + 1
	wsGetNew
	wsGetIntervals
	wsLiveUrl = "https://api-one-wscn.awtmt.com/apiv1/content/lives?channel=global-channel&client=pc%s&limit=20%s&accept=live,vip-live"
)

func (crawler WallStreet) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab.WallStreet",
		Description: "华尔街见闻",
	}
}

// Get 获取 华尔街见闻 快讯
func (crawler WallStreet) Get() error {

	err := crawler.getLiveData(wsGetInit)
	if err != nil {
		logrus.Fatalf("华尔街见闻初始化失败 ：%s", err)
	}
	logrus.Infof("%s 初始化成功",crawler.Config().Description)

	// 定时抓取
	wsCronCrawler()
	return nil
}

func runLatest() {
	_ = WallStreet{}.getLiveData(wsGetNew)
}

func runIntervals() {
	_ = WallStreet{}.getLiveData(wsGetIntervals)
}

func (crawler WallStreet) getLiveData(getType int) error {



	var url string
	switch getType {
	case wsGetInit,wsGetNew:
		logrus.Infof("%s 开始数据采集",crawler.Config().Description)
		url = fmt.Sprintf(wsLiveUrl, "", "&first_page=true")
	case wsGetIntervals:
		// 获取cursor
		logrus.Infof("%s 开始区间数据采集",crawler.Config().Description)
		cursor := service.GetLivesCursor(config.WallStreetLivesCrawler, 0)
		url = fmt.Sprintf(wsLiveUrl, "&cursor="+strconv.FormatInt(cursor, 10), "&first_page=false")
	}

	resp, err := GetUrlData(url,"json")

	if err != nil {
		logrus.Errorf("华尔街见闻数据获取失败 url %s ,error %s ", url, err)
		return err
	}
	logrus.Infof("%s 开始数据解析",crawler.Config().Description)
	lives, err := crawler.respParse(resp)
	if err != nil || len(lives.LivesList) == 0 {
		logrus.Errorf("华尔街见闻数据解析失败 resp %s ,error %s ", resp, err)
		return err
	}

	// new or intervals
	logrus.Infof("%s采集到数据 %d 条， 开始保存数据",crawler.Config().Description,len(lives.LivesList))

	err = service.CreateWallStreetLivesData(lives.LivesList, lives.NextCursor)

	if err != nil {
		logrus.Errorf("华尔街见闻数据保存失败 %s", err)
		return err
	}

	switch getType {
	case wsGetInit:
		logrus.Infof("%s 开始初始化游标 : %d",crawler.Config().Description,lives.NextCursor)

		err = service.InitLivesIntervals(lives.NextCursor, config.WallStreetLivesCrawler, 0)
		if err != nil {
			logrus.Fatalf("华尔街见闻初始化游标失败 %s", err)
		}
	case wsGetIntervals:
		logrus.Infof("%s 开始更新游标 : %d",crawler.Config().Description,lives.NextCursor)
		err = service.UpdateLivesCursor( config.WallStreetLivesCrawler,lives.NextCursor, 0)
		if err != nil {
			logrus.Fatalf("华尔街见闻更新游标失败 %s", err)
		}

	}

	return err
}
func wsCronCrawler() {
	_cron := cron.New()
	err := _cron.AddFunc("@every 1m", runLatest)
	err = _cron.AddFunc("@every 9m", runIntervals)
	if err != nil {
		logrus.Fatalf("华尔街见闻定时器启动失败 %s", err)
	}
	logrus.Infof("华尔街见闻定时采集启动成功 every 1m runLatest every 9m runIntervals")

	_cron.Start()
}

func (crawler WallStreet) respParse(resp string) (service.WallStreetLives, error) {

	var err error
	resJson := gjson.Parse(resp)

	// 解析live
	nextCur := resJson.Get("data.next_cursor").Int()
	respData := resJson.Get("data.items")


	if !respData.Exists() {
		err = errors.New("json parse error :live list no data")
	}
	// slice of live list
	lives := make([]service.WallStreetLivesItems, 0)

	respData.ForEach(func(key, value gjson.Result) bool {

		var items service.WallStreetLivesItems

		value.ForEach(func(k, v gjson.Result) bool {

			switch k.Str {
			case "author":
				items.Author = v.String()
			case "content":
				items.Content = v.Str
			case "content_more":
				items.ContentMore = v.Str
			case "content_text":
				items.ContentText = v.Str
			case "cover_images":
				items.CoverImages = v.String()
			case "display_time":
				items.DisplayTime = v.Int()
			case "id":
				items.Id = v.Int()
			case "images":
				items.Images = v.String()
			case "title":
				items.Title = v.Str
			case "uri":
				items.Uri = v.Str
			}
			return true
		})

		lives = append(lives, items)
		return true
	})

	data := service.WallStreetLives{
		LivesList:  lives,
		NextCursor: nextCur,
	}
	return data, err
}
