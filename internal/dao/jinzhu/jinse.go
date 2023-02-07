/**
 * @Author: scshark
 * @Description:
 * @File:  jinse
 * @Date: 12/12/22 4:29 PM
 */
package jinzhu

import (
	"SecCrawler/internal/core"
	"SecCrawler/internal/model"
	"gorm.io/gorm"
)

var (
	_ core.JinseService = (*jinseServant)(nil)
)

type jinseServant struct {
	db *gorm.DB
}

func newJinseService(db *gorm.DB) core.JinseService {

	return &jinseServant{
		db: db,
	}
}

func (j *jinseServant) CreateMessage(js *model.Jinse, msgData []model.Jinse) (data *model.Jinse, err error) {

	data, err = js.Create(j.db, msgData)

	return data, err
}


func (j *jinseServant) LiveIsExist(liveId int64) (data *model.Jinse, err error) {

	data, err = (&model.Jinse{}).LiveIsExist(j.db,liveId)
	return data, err

}

