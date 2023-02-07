/**
 * @Author: scshark
 * @Description:
 * @File:  dyhjw
 * @Date: 1/6/23 2:28 PM
 */
package core

import "SecCrawler/internal/model"

type DyhjwService interface {
	CreateDyhjwLives(dh *model.Dyhjw,items []model.Dyhjw) (*model.Dyhjw, error)
	DyhjwLivesExist(liveId string) bool
}
