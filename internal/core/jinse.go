/**
 * @Author: scshark
 * @Description:
 * @File:  jinse
 * @Date: 12/12/22 4:34 PM
 */
package core

import "SecCrawler/internal/model"

type JinseService interface {
	CreateMessage(js *model.Jinse,msgData []model.Jinse) (*model.Jinse, error)
	LiveIsExist(liveId int64) (*model.Jinse, error)
}


type JinseResp struct {
	Items []*model.Jinse
	Total int64
}


