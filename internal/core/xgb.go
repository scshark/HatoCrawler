/**
 * @Author: scshark
 * @Description:
 * @File:  xgb
 * @Date: 1/10/23 12:35 PM
 */
package core

import "SecCrawler/internal/model"

type XgbService interface {
	CreateXgbLives(dh *model.Xgb,items []model.Xgb) (*model.Xgb, error)
	XgbLivesIsExist(liveId int64) bool
}