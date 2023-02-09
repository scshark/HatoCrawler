/**
 * @Author: scshark
 * @Description:
 * @File:  xgb
 * @Date: 1/9/23 5:05 PM
 */
package service

import (
	"HatoCrawler/internal/model"
	"github.com/sirupsen/logrus"
)

type XgbLivesCursor int64
type XgbLivesItems struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Image         string `json:"image"`
	LiveCreatedAt int64  `json:"live_created_at"`
	SubjIds       string `json:"subj_ids"`
	Uri           string `json:"uri"`
	Tags          string `json:"tags"`
	OriginaUrl    string `json:"origina_url"`
	Source        string `json:"source"`
}

func CreateXgbLivesData(lives []XgbLivesItems) error {

	var err error
	mData := make([]model.Xgb, 0)

	for _, items := range lives {

		// 已存在
		if ds.XgbLivesIsExist(items.Id) {
			continue
		}

		mData = append(mData, model.Xgb{
			Model:         &model.Model{ID: items.Id},
			Title:         items.Title,
			Summary:       items.Summary,
			Image:         items.Image,
			LiveCreatedAt: items.LiveCreatedAt,
			SubjIds:       items.SubjIds,
			Uri:           items.Uri,
			Tags:          items.Tags,
			OriginaUrl:    items.OriginaUrl,
			Source:        items.Source,
		})
	}
	if len(mData) > 0 {
		_, err = ds.CreateXgbLives(&model.Xgb{}, mData)
		if err != nil {
			logrus.Errorf("XGB数据保存失败 : %s", err)
		}
		logrus.Infof("选股宝 保存数据成功 %d 条",len(mData))

		return err
	}
	return nil
}

