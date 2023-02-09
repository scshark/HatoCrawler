package lab

import (
	"HatoCrawler/config"
	"HatoCrawler/internal/service"
	"HatoCrawler/register"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	// "time"
	// "github.com/mmcdole/gofeed"
)

type Jinse struct{}

const (
	jinseInit = iota + 1
	jinseGetNew
	jinseGetIntervals
	liveUrl = "https://api.jinse.com/noah/v2/lives?limit=20&source=web&flag=down&id=%s&category=0"
)

func (crawler Jinse) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab.Jinse",
		Description: "金色财经",
	}
}

// Get 获取 金色财经快讯 暂时统一入口
func (crawler Jinse) Get() error {

	// 加载数据
	// 初始化信息
	err := crawler.getLiveData(jinseInit)
	if err == nil {
		logrus.Infof("金色财经初始化采集成功")
	}
	// 开启定时刷新数据
	err = jinseCronCrawler()
	return err
}

func jinseCronCrawler() error {

	_cron := cron.New()
	err := _cron.AddFunc("@every 1m", runRefresh)
	err = _cron.AddFunc("@every 7m", runHistory)
	if err != nil {
		logrus.Error("金色财经定时器错误%s", err.Error())
	}
	_cron.Start()
	logrus.Infof("金色财经定时采集已开启 @every 1m runRefresh @every 7m runHistory")
	return err
}

func (crawler Jinse) respParse(resp string) (service.JinseLiveList, error) {

	var err error
	resJson := gjson.Parse(resp)

	topId := resJson.Get("top_id").Int()
	bottomId := resJson.Get("bottom_id").Int()

	// 解析live
	liveList := resJson.Get("list.0.lives")

	if !liveList.Exists() {
		err = errors.New("json parse error :live list no data")
		logrus.Errorf("金色财经数据解析错误 ：%s", err)

	}
	// slice of live list
	s := make([]service.JinseLiveInfo, 0)
	liveList.ForEach(func(key, value gjson.Result) bool {

		data := value.Map()

		// images 需要处理
		var mImages []byte

		if i := data["images"].String(); i != "" {

			mImages, err = json.Marshal(gjson.Parse(i).Value())

			if err != nil {
				logrus.Errorf("金色财经数据解析 json marshal error %s", err)
				return true
			}
		}
		s = append(s, service.JinseLiveInfo{
			LiveId:        value.Get("id").Int(),
			Content:       data["content"].Str,
			Images:        string(mImages),
			ContentPrefix: data["content_prefix"].Str,
			LinkName:      data["link_name"].Str,
			Link:          data["link"].Str,
			LiveCreatedAt: value.Get("live_created_at").Int(),
			CreatedAtZh:   data["created_at_zh"].Str,
		})
		return true
	})

	return service.JinseLiveList{
		TopId:    topId,
		BottomId: bottomId,
		LiveData: s,
	}, err
}

// topId 0 :最新消息   >0 历史消息
func (crawler Jinse) getLiveData(getType uint) error {

	var topId string

	switch getType {
	case jinseInit:
		logrus.Infof("%s 开始初始化采集",crawler.Config().Description)
		topId = "0"
	case jinseGetNew:
		logrus.Infof("%s 开始新数据采集",crawler.Config().Description)
		topId = "0"
	case jinseGetIntervals:
		logrus.Infof("%s 开始区间数据采集",crawler.Config().Description)
		indexId := service.GetLivesCursor(config.JinseLivesCrawler, 0)
		topId = strconv.FormatInt(indexId, 10)
	}

	url := fmt.Sprintf(liveUrl, topId)
	// 请求数据
	respData, err := GetUrlData(url,"json")

	if err != nil {
		logrus.Errorf("金色财经获取数据失败 url %s, error : %s", url, err)
		return err
	}

	// 解析数据
	logrus.Infof("%s 开始数据解析",crawler.Config().Description)
	list, err := crawler.respParse(respData)
	if err != nil {
		logrus.Errorf("金色财经解析数据失败 respData %s, error : %s", respData, err)
		return err
	}

	logrus.Infof("金色财经获取数据 %d 条，开始进行数据保存",len(list.LiveData))
	// 保存列表
	err = service.CreateJinseLiveData(list)
	if err != nil {
		logrus.Errorf("金色财经保存数据失败 error : %s", err)
		return err
	}

	switch getType {
	case wsGetInit:
		logrus.Infof("%s 开始初始化游标 %d",crawler.Config().Description,list.BottomId)

		err = service.InitLivesIntervals(list.BottomId, config.JinseLivesCrawler, 0)
		if err != nil {
			logrus.Fatalf("金色财经初始化游标失败 %s", err)
		}
	case wsGetIntervals:
		logrus.Infof("%s 开始更新游标 %d",crawler.Config().Description,list.BottomId)
		err = service.UpdateLivesCursor(config.JinseLivesCrawler, list.BottomId, 0)
		if err != nil {
			logrus.Fatalf("金色财经更新游标失败 %s", err)
		}

	}

	return err
}

func runRefresh() {
	_ = Jinse{}.getLiveData(jinseGetNew)
}
func runHistory() {
	_ = Jinse{}.getLiveData(jinseGetIntervals)
}
