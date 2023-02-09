// Package service /**
package service

import (
	"HatoCrawler/config"
	"HatoCrawler/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	JinseLiveRedisKey = "jinse_live_index_id"
)

type JinseLiveList struct {
	TopId    int64
	BottomId int64
	LiveData []JinseLiveInfo
}

type JinseLiveInfo struct {
	LiveId        int64  `json:"live_id"`
	Content       string `json:"content"`
	ContentPrefix string `json:"content_prefix"`
	Images        string `json:"images"`
	LinkName      string `json:"link_name"`
	Link          string `json:"link"`
	LiveCreatedAt int64  `json:"live_created_at"`
	CreatedAtZh   string `json:"created_at_zh" `
}

func CreateJinseLiveData(l JinseLiveList) error {
	// 是否已经存在
	//
	mData := make([]model.Jinse, 0)

	isUpdate := true
	for _, item := range l.LiveData {

		// 跳过已有数据
		if _, e := ds.LiveIsExist(item.LiveId); e == nil {
			isUpdate = false
			continue
		}
		mData = append(mData, model.Jinse{
			TopId:         l.TopId,
			BottomId:      l.BottomId,
			LiveId:        item.LiveId,
			Content:       item.Content,
			ContentPrefix: item.ContentPrefix,
			Images:        item.Images,
			Link:          item.Link,
			LinkName:      item.LinkName,
			LiveCreatedAt: item.LiveCreatedAt,
			CreatedAtZh:   item.CreatedAtZh,
		})
	}
	// 无需更新
	if len(mData) == 0 {
		return nil
	}

	// 更新indexId数据
	if isUpdate {
		err := UpdateLivesCursor(config.JinseLivesCrawler, l.BottomId, 0)

		if err != nil {
			logrus.Errorf("金色财经更新cursor失败 error %s ", err)
		}
	}

	logrus.Infof("金色财经更新数据 %d条",len(mData))
	// 保存数据
	_, err := ds.CreateMessage(
		&model.Jinse{},
		mData,
	)

	return err
}
