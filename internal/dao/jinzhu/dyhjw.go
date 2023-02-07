/**
 * @Author: scshark
 * @Description:
 * @File:  dyhjw
 * @Date: 1/6/23 2:27 PM
 */
package jinzhu

import (
	"SecCrawler/internal/core"
	"SecCrawler/internal/model"
	"gorm.io/gorm"
)

var (
	_ core.DyhjwService = (*dyhjwServant)(nil)
)

type dyhjwServant struct {
	db *gorm.DB
}

func newDyhjwService(db *gorm.DB) core.DyhjwService {

	return &dyhjwServant{
		db: db,
	}
}

func (d *dyhjwServant)CreateDyhjwLives(dh *model.Dyhjw,items []model.Dyhjw) (*model.Dyhjw, error)  {
	data, err := dh.Create(d.db, items)
	return data, err
}

func (d *dyhjwServant)DyhjwLivesExist(liveId string) bool  {

	var dh = &model.Dyhjw{
		LiveId:liveId,
	}
	lives,err := dh.First(d.db)

	if err != nil  {
		return false
	}

	if lives.Model != nil && lives.Model.ID > 0 {
		return true
	}
	return false
}