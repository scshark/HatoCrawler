/**
 * @Author: scshark
 * @Description:
 * @File:  WallStreet
 * @Date: 12/29/22 3:19 PM
 */
package service

import (
	"HatoCrawler/internal/model"
	"github.com/sirupsen/logrus"
)


type WallStreetLives struct {
	NextCursor int64
	LivesList []WallStreetLivesItems
}
type WallStreetLivesItems struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Uri         string `json:"uri"`
	DisplayTime int64 `json:"display_time"`
	CoverImages string `json:"cover_images"`
	Content     string `json:"content"`
	ContentText string  `json:"content_text"`
	ContentMore string `json:"content_more" `
	Images      string `json:"images" `
	Author      string `json:"author" `
}


func CreateWallStreetLivesData(lives []WallStreetLivesItems,next int64)error{

	mData := make([]model.WallStreet, 0)

	// var beginIntervals int64
	for _, items := range lives {

		if _,err := ds.WSLivesIsExist(items.Id);err == nil {
			// 记录已经存在
			continue
		}
		mData = append(mData,model.WallStreet{
			Model:         &model.Model{ID: items.Id},
			Title:items.Title,
			ContentText:items.ContentText,
			ContentMore:items.ContentMore,
			Content:items.Content,
			Images:items.Images,
			CoverImages:items.CoverImages,
			Author:items.Author,
			DisplayTime:items.DisplayTime,
			Uri:items.Uri,
		})

		// beginIntervals = items.Id
	}
	if len(mData) > 0{
		_, err := ds.CreateWsLives(&model.WallStreet{}, mData)
		if err != nil {
			logrus.Errorf("华尔街见闻数据保存失败 %s",err)
		}
		logrus.Infof("华尔街见闻数据保存成功 %d 条",len(mData))

		return err
	}

	return nil
}
