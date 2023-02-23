/**
 * @Author: scshark
 * @Description:
 * @File:  wallstreet
 * @Date: 12/29/22 4:37 PM
 */
package model

import "gorm.io/gorm"

type Tabler interface {
	TableName() string
}

type WallStreet struct {
	*Model
	Title       string `json:"title"`
	Uri         string `json:"uri"`
	DisplayTime int64  `json:"display_time"`
	CoverImages string `json:"cover_images"`
	Content     string `json:"content"`
	ContentText string `json:"content_text"`
	ContentMore string `json:"content_more" `
	Images      string `json:"images" `
	Author      string `json:"author" `
}

func (WallStreet)TableName()string  {
	return "ht_wallstreet"
}
func (w *WallStreet) Create (db *gorm.DB,items []WallStreet) (*WallStreet,error){
	err := db.Model(&w).Create(items).Error
	return  w,err
}

func (w *WallStreet) FirstById (db *gorm.DB,id int64)(*WallStreet,error){
	err := db.First(&w, id).Error
	return w,err
}