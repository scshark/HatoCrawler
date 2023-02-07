/**
 * @Author: scshark
 * @Description:
 * @File:  twitter
 * @Date: 1/17/23 3:38 PM
 */
package model

import (
	"gorm.io/gorm"
)


type TwitterUser struct {
	*Model
	TweetUserId      string `json:"tweet_user_id"`
	Name             string `json:"name"`
	ScreenName       string `json:"screen_name"`
	Location         string `json:"location"`
	Description      string `json:"description"`
	Urls             string `json:"urls"`
	DescriptionUrls  string `json:"description_urls"`
	ProfileImageUrl  string `json:"profile_image_url"`
	ProfileBannerUrl string `json:"profile_banner_url"`
	TweetCreatedAt   int64  `json:"tweet_created_at"`
	TwitterLoadTime  int64  `json:"twitter_load_time"`
	LoadOlderTime    int64  `json:"load_older_time"`
	LoadType         int64  `json:"load_type"`
}

func (t *TwitterUser) Create(db *gorm.DB) (*TwitterUser, error) {
	err := db.Create(&t).Error
	return t, err
}

func (t *TwitterUser) First(db *gorm.DB) (*TwitterUser, error) {
	var twUser TwitterUser
	if t.Model != nil && t.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", t.Model.ID, 0)
	} else if t.TweetUserId != "" {
		db = db.Where("tweet_user_id= ? AND is_del = ?", t.TweetUserId, 0)
	} else if t.LoadType > 0 {

		db = db.Where("load_type = ?", 3).Or("load_type = ?", t.LoadType)
		switch t.LoadType {
		case 1:
			db = db.Order("twitter_load_time asc")
		case 2:
			db = db.Order("load_older_time asc")
		}
	}

	err := db.Limit(1).Find(&twUser).Error
	return &twUser, err
}

func (t *TwitterUser) Update(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Omit("id").Save(&t).Error
}
