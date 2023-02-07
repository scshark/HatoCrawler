/**
 * @Author: scshark
 * @Description:
 * @File:  dyhjw
 * @Date: 1/6/23 12:25 PM
 */
package service

import (
	"SecCrawler/internal/model"
)

type DyhjwLivesItems struct {
	Id  string  `json:"id"`
	Content string `json:"content"`
	DisplayTime int64 `json:"display_time"`
}

func CreateDyhjwLivesData(lives []DyhjwLivesItems)error {

	mData := make([]model.Dyhjw, 0)

	for _, items := range lives {

		// 已存在
		if ds.DyhjwLivesExist(items.Id) {
			continue
		}
		
		mData = append(mData,model.Dyhjw{
			LiveId:items.Id,
			Content:items.Content,
			DisplayTime:items.DisplayTime,
			Nonce:items.Id[14:],
			DisplayDatetime:items.Id[:14],
		})
	}
	if len(mData) > 0{
		_, err := ds.CreateDyhjwLives(&model.Dyhjw{}, mData)
		return err
	}

	return nil
}
