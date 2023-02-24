/**
 * @Author: scshark
 * @Description:
 * @File:  jinse
 * @Date: 12/12/22 4:41 PM
 */
package model

import (
	"gorm.io/gorm"
)

type Jinse struct {

	*Model
	TopId int64 `json:"top_id"`
	BottomId int64 `json:"bottom_id"`
	LiveId int64 `json:"live_id"`
	Content string `json:"content"`
	ContentPrefix string `json:"content_prefix"`
	Images string `json:"images"`
	LinkName string `json:"link_name"`
	Link string `json:"link"`
	LiveCreatedAt int64 `json:"live_created_at"`
	CreatedAtZh string `json:"created_at_zh" `
}

func (j *Jinse)Create (db *gorm.DB,msgData []Jinse)(*Jinse,error) {
	err := db.Model(&j).Create(msgData).Error
	return  j,err
}


func(j *Jinse)LiveIsExist(db *gorm.DB,liveID int64)(*Jinse,error) {
	err := db.Where("live_id = ?",liveID).First(&j).Error
	return j,err
}


func (j *Jinse) First(db *gorm.DB,c *ConditionsT) (*Jinse, error) {
	var js Jinse
	if j.Model != nil && j.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", j.Model.ID, 0)
	}
	for k, v := range *c {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}
	err := db.Limit(1).Find(&js).Error
	return &js, err
}