/**
 * @Author: scshark
 * @Description:
 * @File:  xgb
 * @Date: 1/10/23 12:35 PM
 */
package jinzhu

import (
	"HatoCrawler/internal/core"
	"HatoCrawler/internal/model"
	"fmt"
	"gorm.io/gorm"
)

var (
	_ core.XgbService = (*xgbServant)(nil)
)

type xgbServant struct {
	db *gorm.DB
}

func newXgbService(db *gorm.DB) core.XgbService {

	return &xgbServant{
		db: db,
	}
}

func (x *xgbServant)CreateXgbLives(xgb *model.Xgb,items []model.Xgb) (*model.Xgb, error)  {
	data, err := xgb.Create(x.db, items)
	return data, err
}

func (x *xgbServant)XgbLivesIsExist(liveId int64) bool  {

	var xgb = &model.Xgb{
		Model:&model.Model{
			ID:liveId,
		},
	}
	lives,err := xgb.First(x.db)

	if err != nil  {
		fmt.Printf("lives first err %s",err)
		return false
	}

	if lives.Model != nil && lives.Model.ID > 0 {
		return true
	}
	return false
}