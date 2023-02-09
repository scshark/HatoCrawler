/**
 * @Author: scshark
 * @Description:
 * @File:  wallstreet
 * @Date: 12/29/22 5:02 PM
 */
package core

import "HatoCrawler/internal/model"

type WallStreetService interface {
	CreateWsLives(ws *model.WallStreet,items []model.WallStreet) (*model.WallStreet, error)
	WSLivesIsExist(liveId int64) (*model.WallStreet, error)
}