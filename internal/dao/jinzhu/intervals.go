/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 12/30/22 3:56 PM
 */
package jinzhu

import (
	"HatoCrawler/config"
	"HatoCrawler/internal/core"
	"HatoCrawler/internal/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var (
	_ core.IntervalsService = (*intervalsServant)(nil)
)

type intervalsServant struct {
	db *gorm.DB
}

func newIntervalsService(db *gorm.DB) core.IntervalsService {

	return &intervalsServant{
		db: db,
	}
}
func (i *intervalsServant) CreateIntervals(iv *model.Intervals) (*model.Intervals, error) {
	return iv.Create(i.db)
}
func (i *intervalsServant) GetCurrentIntervals(t int64, ex int64) (*model.Intervals, error) {

	c := &model.ConditionsT{
		"is_current = ?": 1,
		"is_completed = ?": 0,
		"type = ?": t,
		"type_extend = ?": ex,
	}
	interval, err := (&model.Intervals{}).First(i.db,c)
	if err != nil {
		return nil, err
	}
	// update current

	if interval.Begin > 0 && interval.Begin < interval.Over {
		interval.IsCurrent = 0
		interval.IsCompleted = 1
		err := interval.Update(i.db)
		if err != nil {
			return nil, err
		}
		// pick 一个新的interval
		return i.PickNewIntervals(t,ex)
	} else {
		return interval, err
	}

}
func (i *intervalsServant) GetIntervalsOverCursor(t int64, ex int64) (overCursor int64,err error) {

	switch t {
	// 获取最大id作为本次区间over cursor
	case config.JinseLivesCrawler:
		c := &model.ConditionsT{
			"created_on < ?": time.Now().Unix() - 35,
			"Order": "top_id DESC ,live_id DESC",
		}
		live,err := (&model.Jinse{}).First(i.db,c)
		if err != nil {
			return 0,err
		}
		overCursor = live.TopId

	case config.WallStreetLivesCrawler:
		c := &model.ConditionsT{
			"created_on < ?": time.Now().Unix() - 35,
			"Order": "display_time DESC",
		}
		live,err := (&model.WallStreet{}).First(i.db,c)
		if err != nil {
			return 0,err
		}
		overCursor = live.DisplayTime
	case config.DyhjwLivesCrawler:
		c := &model.ConditionsT{
			"created_on < ?": time.Now().Unix() - 35,
			"Order": "display_time DESC",
		}
		live,err := (&model.Dyhjw{}).First(i.db,c)
		if err != nil {
			return 0,err
		}
		overCursor,_ = strconv.ParseInt(live.LiveId,10,64)

	case config.XgbLivesCrawler:

		c := &model.ConditionsT{
			"created_on < ?": time.Now().Unix() - 35,
			"Order": "live_created_at DESC",
		}
		live,err := (&model.Xgb{}).First(i.db,c)
		if err != nil {
			return 0,err
		}
		overCursorStr := strconv.FormatInt(live.LiveCreatedAt,10) + "133"

		overCursor,_ = strconv.ParseInt(overCursorStr,10,64)
	}


	return overCursor,nil
}
func (i *intervalsServant) CancelIntervalsCurrent(t int64, ex int64) error {
	var iv = &model.Intervals{
		IsCurrent:  0,
		Type:       t,
		TypeExtend: ex,
	}
	return iv.Updates(i.db, "IsCurrent")
}
func (i *intervalsServant) PickNewIntervals(t int64, ex int64) (*model.Intervals, error) {

	c := &model.ConditionsT{
		"is_current = ?": 0,
		"is_completed = ?": 0,
		"type = ?": t,
		"type_extend = ?": ex,
	}
	interval, err := (&model.Intervals{}).First(i.db,c)
	if err != nil {
		return nil, err
	}
	interval.IsCurrent = 1
	err = interval.Update(i.db)

	return interval, err
}
func (i *intervalsServant) UpdateIntervals(iv *model.Intervals) error {
	return iv.Update(i.db)
}
