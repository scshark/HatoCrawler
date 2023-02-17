/**
 * @Author: scshark
 * @Description:
 * @File:  Twitter
 * @Date: 1/12/23 1:04 PM
 */
package lab

import (
	"HatoCrawler/internal/conf"
	"HatoCrawler/internal/service"
	"HatoCrawler/register"
	"HatoCrawler/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"math/big"
	"net/http"
	"time"
)

type Twitter struct{
	ScreenName []string
}

const (
	twitterGetUrl = "https://socialbearing.com/scripts/get-tweets.php?%s&search=%s&searchtype=user"
)


func (crawler Twitter) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "Lab.Twitter",
		Description: "Twitter",
	}
}

func (crawler Twitter) Get() error {

	var err error

	logrus.Infof("推特采集初始化 列表ScreenName %v ",crawler.ScreenName)
	for _, t := range crawler.ScreenName {
		conf.Redis.SAdd(conf.Ctx,"twitter_init_screenName",t)
	}
	err = crawler.initTwitterLives()
	if err != nil {
		logrus.Errorf("推特采集初始化错误 initTwitterLives err : %s", err)
	}
	// 定时采集
	cronTweetCrawler()
	return err
}

func cronTweetCrawler() {
	_cron := cron.New()
	err := _cron.AddFunc("@daily",runInitLoad)// 每天初始化
	err = _cron.AddFunc("@every 15s",runLoadNewer)// 每15秒更新
	err = _cron.AddFunc("@every 5m",runLoadOlder)// 每5分钟获取一次旧消息
	if err != nil {
		logrus.Fatalf("推特采集定时器启动失败 error: %s\n", err.Error())
	}
	logrus.Infof("Twitter 定时采集已启动 daily runInitLoad - every 1mrunLoadNewer - every 5m runLoadOlder")
	_cron.Start()
}

func runInitLoad(){
	err := Twitter{}.initTwitterLives()
	if err != nil {
		logrus.Errorf("推特采集初始化错误 initTwitterLives err : %s", err)
	}
}
func runLoadNewer() {
	err := Twitter{}.loadNewerTwitter()
	if err != nil {
		logrus.Errorf("推特新消息获取失败 loadNewerTwitter err : %s", err)
	}
}
func runLoadOlder() {
	err := Twitter{}.loadOlderTwitter()
	if err != nil {
		logrus.Errorf("推特旧消息获取失败 loadOlderTwitter err : %s", err)
	}
}



/**
 * 加载推特用户最新消息
 * @author: sc.shark
 * @date: Created in 2023-02-03 14:40:20
 * @Description:
 * @receiver crawler
 * @return error
 */
func (crawler Twitter) loadNewerTwitter() error {

	logrus.Info("推特开始采集新信息 ")
	// 获取用户的 name
	tweetUser := service.GetTweetUserForLoad(service.LoadNewer)

	logrus.Infof("此次采集推特新消息的 用户 ：%s ",tweetUser.Name)

	if tweetUser == nil {
		return errors.New("GetTweetUserForLoad error : record not found")
	}

	// 获取用户最新的推特id

	lastTw := service.GetTwUserLastTweet(tweetUser.ID)

	logrus.Infof("此次采集推特新消息的 用户 最近一条推特 id ：%s ",lastTw.IdStr)

	getUrl := ""
	if lastTw.IdStr == ""{
		getUrl = fmt.Sprintf(twitterGetUrl, "sid=0", tweetUser.ScreenName)
	}else{
		// 获取信息地址
		getUrl = fmt.Sprintf(twitterGetUrl, "sid="+lastTw.IdStr, tweetUser.ScreenName)
	}

	logrus.Infof("此次采集推特新消息的 用户 ：%s ,获取链接 %s ",tweetUser.Name,getUrl)

	resp, err := crawler.GetData(getUrl)
	if err != nil {
		logrus.Warnf("loadNewerTwitter  GetData error:%s",err)
		return err
	}

	logrus.Infof("此次采集推特新消息的 用户 ：%s , 开始解析数据",tweetUser.Name)

	tw, err := crawler.respParse(resp)

	if err != nil {
		logrus.Warnf("loadNewerTwitter  respParse error: %s",err)
		return err
	}

	if len(tw.Items) == 0  {
		logrus.Infof("推特用户%s 数据已经最新，更新用户推特更新时间LoadTime",tweetUser.Name)
		err = service.UpdateUserLoadTime(service.LoadNewer,tweetUser)
		if err !=nil {
			logrus.Errorf("更新推特用户更新时间失败 err :%s",err)
			return  err
		}
		return nil
	}

	// save user
	if tw.User.IdStr == "" {
		logrus.Warnf("当前推特用户ID丢失 err :%s",err)
		return err
	}

	logrus.Infof("推特用户%s 采集到 %d 个用户信息",tweetUser.Name,len(tw.ReplyUser) + 1)

	user := []service.TwitterUser{tw.User}
	for _, u := range tw.ReplyUser {
		user = append(user, u)
	}

	tmpUser, err := service.SaveTweetUser(user)


	if err != nil {
		logrus.Errorf("推特用户保存失败 err :%s",err)
		return err
	}
	logrus.Infof("推特用户%s 用户信息更新成功，共保存 %d 个用户信息",tweetUser.Name,len(tmpUser))


	// save twitter
	logrus.Infof("推特用户%s，用户推特id %d，采集类型：新消息， 此次共采集到 %d 新消息 ，开始保存数据",tweetUser.Name,tweetUser.ID,len(tw.Items))

	err = service.CreateTwitterListData(tweetUser.ID, tw.Items)
	if err != nil {
		logrus.Errorf("推特信息保存失败 err :%s",err)
		return err
	}
	// save reply twitter
	logrus.Infof("推特用户%s，采集类型：新消息（回复）， 此次共采集到 %d 回复消息 ，开始保存数据",tweetUser.Name,len(tw.ReplyItems))
	err = service.CreateTwitterListData(0, tw.ReplyItems)
	if err != nil {
		logrus.Errorf("推特回复信息保存失败 err :%s",err)
		return err
	}
	logrus.Infof("推特用户%s，更新用户推特采集新消息时间 LoadTime",tweetUser.Name)
	// 更新用户获取信息时间
	err = service.UpdateUserLoadTime(service.LoadNewer,tweetUser)
	if err !=nil {
		logrus.Errorf("更新推特用户更新时间失败 err :%s",err)
		return  err
	}
	logrus.Infof("推特用户%s 完成新信息采集",tweetUser.Name)

	return err
}



/**
 * 获取推特用户旧消息
 * @author: sc.shark
 * @date: Created in 2023-02-03 16:24:13
 * @Description:
 * @receiver crawler
 * @return error
 */
func (crawler Twitter) loadOlderTwitter() error {


	logrus.Info("推特开始采集区间信息 ")

	// 获取用户的 name
	tweetUser := service.GetTweetUserForLoad(service.LoadOlder)

	logrus.Infof("推特此次获取区间信息用户为 ： %s ",tweetUser.Name)

	if tweetUser == nil {
		return errors.New("GetTweetUserForLoad error : record not found")
	}

	lastTw := service.GetTwUserFirstTweet(tweetUser.ID)
	logrus.Infof("推特此次获取区间信息用户最新一条推特ID为 ： %s ",lastTw.IdStr)


	getUrl := fmt.Sprintf(twitterGetUrl, "oid="+lastTw.IdStr, tweetUser.ScreenName)


	logrus.Infof("推特此次获取区间信息用户为 ： %s , 获取链接 %s",tweetUser.Name,getUrl)
	resp, err := crawler.GetData(getUrl)

	if err != nil {
		logrus.Errorf("推特信息获取失败 respParse error: %s",err)
		return err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 开始解析数据",tweetUser.Name)

	tw, err := crawler.respParse(resp)

	if err != nil {
		logrus.Errorf("推特信息解析失败 respParse error: %s",err)
		return err
	}

	if len(tw.Items) == 0  {
		logrus.Infof("推特此次获取区间信息用户为 ： %s , 没有获取到新信息，更新用户区间采集信息时间 load_older_time",tweetUser.Name)
		err = service.UpdateUserLoadTime(service.LoadOlder,tweetUser)
		if err !=nil {
			logrus.Errorf("更新推特用户更新时间失败 err :%s",err)
			return  err
		}
		return nil
	}

	// save user
	if tw.User.IdStr == "" {
		logrus.Errorf("当前推特用户ID丢失 err :%s",err)
		return err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 共获取 %d 条用户数据",tweetUser.Name,len(tw.ReplyUser) + 1)

	user := []service.TwitterUser{tw.User}
	for _, u := range tw.ReplyUser {
		user = append(user, u)
	}


	tmpUser, err := service.SaveTweetUser(user)

	if err != nil {
		logrus.Errorf("推特用户信息保存失败 err :%s",err)
		return err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 共保存 %d 条用户数据",tweetUser.Name,len(tmpUser))

	logrus.Infof("推特此次获取区间信息用户为 ： %s , 共采集到 %d 条区间信息",tweetUser.Name,len(tw.Items))
	// save twitter
	err = service.CreateTwitterListData(tweetUser.ID, tw.Items)
	if err != nil {
		logrus.Errorf("推特信息保存失败 err :%s",err)
		return err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 共采集到 %d 条区间回复信息",tweetUser.Name,len(tw.ReplyItems))
	// save reply twitter
	err = service.CreateTwitterListData(0, tw.ReplyItems)
	if err != nil {
		logrus.Errorf("推特回复信息保存失败 err :%s",err)
		return err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 更新用户区间采集时间 load older time",tweetUser.Name)
	err = service.UpdateUserLoadTime(service.LoadOlder,tweetUser)
	if err !=nil {
		logrus.Errorf("推特用户信息更新时间保存 err :%s",err)
		return  err
	}
	logrus.Infof("推特此次获取区间信息用户为 ： %s , 区间数据采集完成",tweetUser.Name)

	return err

}

/**
 * 初始化推特信息
 * @author: sc.shark
 * @date: Created in 2023-02-03 16:12:03
 * @Description:
 * @receiver crawler
 * @return error
 */
func (crawler Twitter) initTwitterLives() error {

	var err error
	// 需要初始化推特用户的列表

	initTweetUserList := conf.Redis.SMembers(conf.Ctx,"twitter_init_screenName").Val()
	logrus.Infof(" %s 初始化开始，初始化用户列表 %v",crawler.Config().Description,initTweetUserList)
	for _, t := range initTweetUserList {

		getUrl := fmt.Sprintf(twitterGetUrl, "sid=0", t)

		logrus.Infof("初始化获取用户：%s，获取推特信息地址 : %s",t,getUrl)
		resp, err := crawler.GetData(getUrl)

		if err != nil {
			return err
		}

		logrus.Infof("初始化获取用户：%s，开始解析数据",t)
		tw, err := crawler.respParse(resp)

		if err != nil {
			return err
		}
		if len(tw.Items) == 0  {
			logrus.Warnf("推特用户%s 没有获取到数据",t)
			continue
		}
		logrus.Infof("已获取到用户推特信息: %s，推特 %d 条，回复用户 %d 个 ，回复推特 %d 条",t,len(tw.Items),len(tw.ReplyUser),len(tw.ReplyItems))

		// save user
		if tw.User.IdStr == "" {
			logrus.Errorf("当前推特用户ID丢失 err :%s",err)
			return err
		}
		logrus.Infof("初始化获取用户：%s, 开始保存用户信息, 需要保存 %d 个用户",t,len(tw.ReplyUser) + 1)
		user := []service.TwitterUser{tw.User}
		for _, u := range tw.ReplyUser {
			user = append(user, u)
		}
		sTwUser, err := service.SaveTweetUser(user)

		logrus.Infof("初始化获取用户：%s，此次保存 %d 用户",t,len(sTwUser))

		if err != nil {
			return err
		}

		if sTwUser[0].UserId == 0 {
			return errors.New("推特用户id错误")
		}

		tweetUser := sTwUser[0]
		// save twitter
		logrus.Infof("初始化获取用户：%s, 开始保存用户推特, 需要保存 %d 条推特",t,len(tw.Items))

		err = service.CreateTwitterListData(tweetUser.UserId, tw.Items)
		if err != nil {
			return err
		}
		// save reply twitter
		logrus.Infof("初始化获取用户：%s, 开始保存用户回复推特, 需要保存 %d 条回复推特",t,len(tw.ReplyItems))
		err = service.CreateTwitterListData(0, tw.ReplyItems)
		if err != nil {
			return err
		}
		// sleep 30s
		logrus.Infof("推特用户%s 完成初始化信息采集",t)

		logrus.Info("推特采集初始化 time sleep 55 second")
		time.Sleep(55 * time.Second)

	}
	return err

}

// 请求获取数据
func (crawler Twitter) GetData(url string) (string, error) {

	client := utils.CrawlerClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	w, err := io.Copy(&buffer, resp.Body)

	if err != nil || w < 99999 {
		return "", err
	}

	return buffer.String(), err
}

// 解析数据
func (crawler Twitter) respParse(resp string) (service.TwitterParse, error) {

	var err error
	var user service.TwitterUser

	resJson := gjson.Parse(resp)

	// 解析live
	if !resJson.Exists() {
		err = errors.New("没有数据 resp:" + resp)
		return service.TwitterParse{},nil
	}
	// slice of live list
	repUser := make([]service.TwitterUser, 0)
	tw := make([]service.TwitterItems, 0)
	replyTw := make([]service.TwitterItems, 0)
	rpUserId := make([]string, 0)
	// replyUser := make([]service.TwitterUser,0)
	resJson.ForEach(func(key, value gjson.Result) bool {

		if !value.Exists() {
			logrus.Errorf("推特信息解析错误")
			return true
		}

		// 抽出user处理
		// 只保存第一个items的user
		if user.IdStr == "" {
			tUser, err := crawler.tweetUserEntitiesParse(value.Get("user"))
			tUser.LoadType = 3
			if err != nil {
				logrus.Errorf("推特用户主体信息解析错误 %s",err )
			} else {
				user = tUser
			}
		}

		// 解析twitter 主体信息
		tweet, err := crawler.tweetEntitiesParse(value)

		if err != nil {
			logrus.Errorf("解析twitter 主体信息错误 %s",err )
			return true
		}

		tw = append(tw, tweet)

		// 抽取 retweet user
		if value.Get("retweeted_status").Exists() {
			// rpUser是否已经存在slice里

			// 把reply的user 加入到slice
			rpUser, err := crawler.tweetUserEntitiesParse(value.Get("retweeted_status.user"))

			if err != nil {
				logrus.Errorf("json parse retweeted_status tweetUserEntitiesParse error %s",err )
			} else {
				isAppendRpU := false
				rpUser.LoadType = 0

				// 只添加粉丝15W以上的账户
				if rpUser.FollowersCount > 155555 {
					isAppendRpU = true
				}

				for _, i := range rpUserId {
					if i == rpUser.IdStr {
						isAppendRpU = false
					}
				}

				if isAppendRpU{
					rpUserId = append(rpUserId, rpUser.IdStr)
					repUser = append(repUser, rpUser)
				}
			}

			// rp tweet
			reTweet, err := crawler.tweetEntitiesParse(value.Get("retweeted_status"))
			if err != nil {

				logrus.Warnf("tweet parse retweeted_status error %s", err)
				return true
			}

			replyTw = append(replyTw, reTweet)

		}

		return true
	})
	return service.TwitterParse{
		Items:      tw,
		ReplyItems: replyTw,
		User:       user,
		ReplyUser:  repUser,
	}, err
}

// 解析推特主体数据
func (crawler Twitter) tweetEntitiesParse(result gjson.Result) (items service.TwitterItems, err error) {

	reply := map[string]string{}
	result.ForEach(func(k, v gjson.Result) bool {

		switch k.Str {
		case "created_at":
			// Wed Jan 11 14:05:50 +0000 2023
			// 需要处理

			createdAt, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", v.Str)
			if err != nil {
				logrus.Warnf("tw时间解析错误 %s",err)
				return true
			}
			items.CreatedAt = createdAt.Unix()
		case "id":
			items.Id = big.NewInt(v.Int())
		case "id_str":
			items.IdStr = v.Str
		case "full_text":
			items.FullText = v.Str
		case "user":
			items.UserId = big.NewInt(v.Get("id").Int())
		case "entities":
			// hashtags
			// user_mentions
			// urls

			if v.Get("user_mentions.#").Int() > 0 {

				userMen := make([]map[string]string, 0)

				v.Get("user_mentions").ForEach(func(kk, vv gjson.Result) bool {

					// var u map[string]string

					um := vv.Map()
					u := map[string]string{
						"screen_name": um["screen_name"].Str,
						"name":        um["name"].Str,
						"id_str":      um["id_str"].Str,
					}
					userMen = append(userMen, u)
					return true
				})
				if len(userMen) > 0 {
					umJson, err := json.Marshal(userMen)
					if err != nil {
						logrus.Warnf("user_mentions json Marshal error %s",err)
						break
					}
					items.UserMentions = fmt.Sprintf("%s", umJson)
				}

			}
			items.Hashtags = v.Get("hashtags.#.text").String()
			items.Urls = v.Get("urls.#.expanded_url").String()
		case "in_reply_to_status_id_str", "in_reply_to_user_id_str", "in_reply_to_screen_name":
			// 处理回复
			if v.Str == "" {
				break
			}

			reply[k.Str] = v.Str

		case "extended_entities":
			// extended_entities处理 媒体信息
			if media := v.Get("media").Array(); len(media) > 0 {

				tweetMedia := make([]map[string]string, 0)

				for _, v := range media {
					md := v.Map()
					var videoUrl string
					if md["type"].Str == "video" {
						videoUrl = v.Get("video_info.variants.0.url").Str
					}
					tweetMedia = append(tweetMedia, map[string]string{
						"media_url": md["media_url_https"].Str,
						"type":      md["type"].Str,
						"video_url": videoUrl,
					})

				}

				if len(tweetMedia) > 0 {
					tm, err := json.Marshal(tweetMedia)
					if err != nil {
						logrus.Warnf("resp parse extended_entities tweetMedia error  %s",err)
						break
					}
					items.ExtendedEntities = fmt.Sprintf("%s", tm)
				}
			}

		}
		return true
	})
	replyJson, err := json.Marshal(reply)
	if err != nil {
		logrus.Warnf("reply json Marshal  error %s",err)
	} else {
		items.InReplyInfo = fmt.Sprintf("%s", replyJson)
	}
	return items, err
}

// 解析推特用户数据
func (crawler Twitter) tweetUserEntitiesParse(result gjson.Result) (user service.TwitterUser, err error) {

	result.ForEach(func(k, v gjson.Result) bool {

		switch k.Str {
		case "id":
			user.Id = big.NewInt(v.Int())
		case "id_str":
			user.IdStr = v.Str
		case "name":
			user.Name = v.Str
		case "screen_name":
			user.ScreenName = v.Str
		case "location":
			user.Location = v.Str
		case "description":
			user.Description = v.Str
		case "followers_count":
			user.FollowersCount = v.Int()
		case "friends_count":
			user.FriendsCount = v.Int()
		case "profile_image_url_https":
			user.ProfileImageUrl = v.Str
		case "profile_banner_url":
			user.ProfileBannerUrl = v.Str
		case "created_at":
			createdAt, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", v.Str)
			if err != nil {
				logrus.Warnf(" 推特用户信息 respParse CreatedAt time.Parse err %s", err)
				return true
			}

			user.CreatedAt = createdAt.Unix()
		case "entities":

			user.Url = v.Get("url.urls.#.expanded_url").String()
		}
		return true
	})
	return user, err
}

