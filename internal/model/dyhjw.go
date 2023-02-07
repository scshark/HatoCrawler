/**
 * @Author: scshark
 * @Description:
 * @File:  dyhjw
 * @Date: 1/6/23 12:19 PM
 */
package model

import (
	"gorm.io/gorm"
)

type Dyhjw struct {
	*Model
	LiveId          string `json:"live_id"`
	Content         string `json:"content"`
	IsTweet         int64  `json:"is_tweet"`
	DisplayTime     int64  `json:"display_time"`
	DisplayDatetime string  `json:"display_datetime"`
	Nonce           string  `json:"nonce"`
}
func (d *Dyhjw) Create (db *gorm.DB,items []Dyhjw) (*Dyhjw,error){
	err := db.Model(&d).Create(items).Error
	return  d,err
}

func (d *Dyhjw) First(db *gorm.DB) (*Dyhjw, error) {
	var dh Dyhjw
	if d.Model != nil && d.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", d.Model.ID, 0)
	}else if d.LiveId != "" {
		db = db.Where("live_id= ? AND is_del = ?", d.LiveId, 0)
	}else if d.DisplayTime > 0 {
		db = db.Where("display_time= ? AND is_del = ?", d.DisplayTime, 0)
	}
	err := db.Limit(1).Find(&dh).Error
	return &dh, err
}

func (d *Dyhjw) Update(db *gorm.DB) error{
	return db.Model(&Intervals{}).Where("id = ? AND is_del = ?", d.Model.ID, 0).Omit("id").Save(d).Error
}
