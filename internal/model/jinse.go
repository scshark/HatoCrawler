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