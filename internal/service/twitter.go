/**
 * @Author: scshark
 * @Description:
 * @File:  twitter
 * @Date: 1/12/23 7:49 PM
 */
package service

import (
	"HatoCrawler/internal/model"
	"errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

const (
	LoadNewer = iota
	LoadOlder
)
const (
	NoLoad = iota
	OnlyLoadNewer
	OnlyLoadOlder
	LoadAll
)
type TwitterUser struct {
	Id               *big.Int `json:"id"`
	UserId           int64    `json:"user_id"`
	IdStr            string   `json:"id_str"`
	Name             string   `json:"name"`
	ScreenName       string   `json:"screen_name"`
	Location         string   `json:"location"`
	Description      string   `json:"description"`
	Url              string   `json:"url"`
	DescriptionUrls  string   `json:"description_urls"`
	ProfileImageUrl  string   `json:"profile_image_url"`
	ProfileBannerUrl string   `json:"profile_banner_url"`
	CreatedAt        int64    `json:"created_at"`
	FollowersCount        int64    `json:"followers_count"`
	FriendsCount        int64    `json:"friends_count"`
	LoadType        int64    `json:"load_type"`
}

type TwitterItems struct {
	Id               *big.Int `json:"id"`
	UserId           *big.Int `json:"user_id"`
	IdStr            string   `json:"id_str"`
	FullText         string   `json:"full_text"`
	Hashtags         string   `json:"hashtags"`
	UserMentions     string   `json:"user_mentions"`
	Urls             string   `json:"urls"`
	ExtendedEntities string   `json:"extended_entities"`
	// RetweetedItem   ReplyTwitterItems
	InReplyInfo string `json:"in_reply_info"`
	CreatedAt   int64  `json:"created_at"`
}

type TwitterParse struct {
	Items      []TwitterItems
	ReplyItems []TwitterItems
	User       TwitterUser
	ReplyUser  []TwitterUser
}

func GetTweetUserByName(name string) *model.TwitterUser {
	return ds.GetTweetUserByScreenName(name)
}
func SaveTwitterUserItems(u TwitterUser) (twUser *model.TwitterUser, err error) {
	// save user,
	twUser = ds.GetTweetUserByTweetId(u.IdStr)

	if twUser.Model != nil && twUser.Model.ID > 0 {

		tmpUser := twUser

		twUser.Name = u.Name
		twUser.Description = u.Description
		twUser.Location = u.Location
		twUser.DescriptionUrls = u.DescriptionUrls
		twUser.ProfileBannerUrl = u.ProfileBannerUrl
		twUser.ProfileImageUrl = u.ProfileImageUrl
		twUser.TweetCreatedAt = u.CreatedAt
		twUser.Urls = u.Url
		twUser.FriendsCount = u.FriendsCount
		twUser.FollowersCount = u.FollowersCount
		twUser.NeedHatoUpdate = 1
		//twUser.LoadType = u.LoadType

		if tmpUser != twUser {
			logrus.Infof("name: %s ????????????????????????",u.Name)
			err = ds.UpdateTweetUser(twUser)
		}else {
			logrus.Infof("name: %s ???????????????????????????????????????",u.Name)
		}


	} else {
		logrus.Infof("name: %s ????????????????????????",u.Name)
		var user = &model.TwitterUser{
			TweetUserId:      u.IdStr,
			Name:             u.Name,
			ScreenName:       u.ScreenName,
			Description:      u.Description,
			DescriptionUrls:  u.DescriptionUrls,
			Location:         u.Location,
			ProfileBannerUrl: u.ProfileBannerUrl,
			ProfileImageUrl:  u.ProfileImageUrl,
			TweetCreatedAt:   u.CreatedAt,
			LoadType:         3,
			//TwitterLoadTime:  time.Now().Unix(),
			FriendsCount: u.FriendsCount,
			FollowersCount: u.FollowersCount,
			NeedHatoUpdate: 1,
		}
		twUser, err = ds.CreateTweetUser(user)
		if err != nil {
			logrus.Errorf("CreateTweetUser error : %s", err)
		}
	}

 	return twUser, err
}

func CreateTwitterListData(userId int64, tw []TwitterItems) error {

	var err error

	if len(tw) == 0 {
		err = errors.New("??????????????????????????????")
		return err
	}
	mData := make([]model.Twitter, 0)

	for _, items := range tw {

		// ?????????
		if ds.TweetIsExists(items.Id) {
			continue
		}

		if userId == 0 {
			user := ds.GetTweetUserByTweetId(items.UserId.String())
			if user.Model == nil {
				continue
			}
			userId = user.ID
		}

		//if items.UserMentions != "" {
		//	logrus.Warnf("user ------ UserMentions %s", items.UserMentions)
		//}
		mData = append(mData, model.Twitter{
			Model:            &model.Model{ID: items.Id.Int64()},
			IdStr:            items.IdStr,
			HtTwitterUserId:  userId,
			TwitterUserId:    items.UserId.String(),
			FullText:         items.FullText,
			Hashtags:         items.Hashtags,
			UserMentions:     items.UserMentions,
			Urls:             items.Urls,
			ExtendedEntities: items.ExtendedEntities,
			InReplyInfo:      items.InReplyInfo,
			TwCreatedAt:      items.CreatedAt,
		})
	}
	if len(mData) > 0 {
		err = ds.CreateTweet(&model.Twitter{}, mData)
		if err != nil {
			logrus.Errorf("CreateTwitter error : %s", err)
		}
		logrus.Infof("?????????????????????????????? %d ???",len(mData))

	}
	return err
}
func GetTwUserLastTweet(userId int64) *model.Twitter {

	tw, err := ds.GetTweetByUserId(userId, &model.ConditionsT{
		"ORDER": "id desc",
	})
	if err != nil {
		logrus.Warnf("GetLastTweet err %s", err)
	}
	return tw
}

func GetTwUserFirstTweet(userId int64) *model.Twitter {

	tw, err := ds.GetTweetByUserId(userId, &model.ConditionsT{
		"ORDER": "id asc",
	})
	if err != nil {
		logrus.Warnf("GetFirstTweet err %s", err)
	}
	return tw
}

func GetTweetUserForLoad(loadType uint) *model.TwitterUser {

	var tUser = &model.TwitterUser{}
	switch loadType {

	case LoadNewer:
		// ?????????100W???????????????
		tUser = &model.TwitterUser{
			LoadType:        OnlyLoadNewer,
		}
	case LoadOlder:
		tUser = &model.TwitterUser{
			LoadType:      OnlyLoadOlder,
		}
	default:
		return nil
	}
	twUser := ds.GetTweetUserForLoad(tUser)

	// ???????????????????????? 5000 ????????? ??????????????? ????????? OnlyLoadNewer
	twCount, err := ds.CountTweetByUserId(twUser.ID, &model.ConditionsT{
		"ORDER": "id desc",
	})
	if err !=nil {
		logrus.Errorf("CountTweetByUserId ??????????????????????????????????????? ?????? %s",err)
	}
	logrus.Infof("CountTweetByUserId ??????????????????????????????????????? ?????? %s ,????????? ?????? %d ???",twUser.ScreenName,twCount)
	if twCount > 5000 {
		twUser.LoadType = OnlyLoadNewer
		err = ds.UpdateTweetUser(twUser)

		if err != nil{
			logrus.Errorf("UpdateTweetUser ?????????????????????????????????????????????????????? ?????? %s",err)
		}
	}
	return twUser
}

func UpdateUserLoadTime(loadType uint, user *model.TwitterUser) error {

	switch loadType {
	case LoadNewer:
		user.TwitterLoadTime = time.Now().Unix()
	case LoadOlder:
		user.LoadOlderTime = time.Now().Unix()
	default:
		return errors.New("??????????????????")
	}
	err := ds.UpdateTweetUser(user)

	return err
}

// ????????????????????????
func SaveTweetUser(twUser []TwitterUser) ([]TwitterUser, error) {
	// save user and reply user

	user := make([]TwitterUser,0)
	var err error
	for k, u := range twUser {

		tUser, err := SaveTwitterUserItems(u)

		if err != nil {
			logrus.Errorf(" user SaveTwitterUserItems error: %s\n", err)
			if k == 0 {
				return nil, err
			}
		}
		// ????????????????????????????????????reply user
		if k == 0 && tUser.Model == nil {
			return nil, err
		}
		u.UserId = tUser.ID
		user = append(user,u)
	}

	return user, err
}
