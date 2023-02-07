/**
 * @Author: scshark
 * @Description:
 * @File:  wallstreet
 * @Date: 12/29/22 5:03 PM
 */
package jinzhu

import (
	"SecCrawler/internal/core"
	"SecCrawler/internal/model"
	"gorm.io/gorm"
)

var (
	_ core.WallStreetService = (*wallStreetServant)(nil)
)

type wallStreetServant struct {
	db *gorm.DB
}

func newWallStreetService(db *gorm.DB) core.WallStreetService {

	return &wallStreetServant{
		db: db,
	}
}

func (w *wallStreetServant)CreateWsLives(ws *model.WallStreet,items []model.WallStreet) (*model.WallStreet, error)  {

	data, err := ws.Create(w.db, items)

	return data, err
}

func (w *wallStreetServant)WSLivesIsExist(liveId int64) (*model.WallStreet, error)  {

	data, err := (&model.WallStreet{}).FirstById(w.db,liveId)
	return data, err

}
