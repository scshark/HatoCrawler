/**
 * @Author: scshark
 * @Description:
 * @File:  twitter
 * @Date: 1/17/23 6:58 PM
 */
package model

import (
	"gorm.io/gorm"
)

type Twitter struct {
	*Model
	HtTwitterUserId  int64 `json:"ht_twitter_user_id"`
	TwitterUserId    string   `json:"twitter_user_id"`
	IdStr            string   `json:"id_str"`
	FullText         string   `json:"full_text"`
	Hashtags         string   `json:"hashtags"`
	UserMentions     string   `json:"user_mentions"`
	Urls             string   `json:"urls"`
	ExtendedEntities string   `json:"extended_entities"`
	InReplyInfo      string   `json:"in_reply_info"`
	TwCreatedAt        int64    `json:"tw_created_at"`
}

func (t *Twitter) Create(db *gorm.DB,data []Twitter) (*Twitter, error) {
	err := db.Create(&data).Error
	return t, err
}

func (t *Twitter) First(db *gorm.DB,c *ConditionsT) (*Twitter, error) {
	var tw Twitter
	if t.Model != nil && t.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", t.Model.ID, 0)
	} else if t.IdStr != "" {
		db = db.Where("id_str= ? AND is_del = ?", t.IdStr, 0)
	}else if t.HtTwitterUserId > 0 {
		db = db.Where("ht_twitter_user_id= ? AND is_del = ?", t.HtTwitterUserId, 0)
	}
	for k, v := range *c {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}
	err := db.Limit(1).Find(&tw).Error
	return &tw, err
}

func (t *Twitter) Exists(db *gorm.DB) error {
	var tw Twitter
	if  t.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", t.ID, 0)
	} else if t.IdStr != "" {
		db = db.Where("id_str= ? AND is_del = ?", t.IdStr, 0)
	}
	err := db.First(&tw).Error
	return err
}


func (t *Twitter) Count(db *gorm.DB,c *ConditionsT) (int64,error) {
	var count int64
	if  t.HtTwitterUserId > 0 {
		db = db.Where("ht_twitter_user_id= ? AND is_del = ?", t.HtTwitterUserId, 0)
	}
	for k, v := range *c {
		if k != "ORDER" {
			db = db.Where(k, v)
		}
	}
	if err := db.Model(t).Count(&count).Error; err != nil {
		return 0, err
	}
	return count,nil
}
