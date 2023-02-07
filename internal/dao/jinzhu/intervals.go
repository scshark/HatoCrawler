/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 12/30/22 3:56 PM
 */
package jinzhu

import (
	"SecCrawler/internal/core"
	"SecCrawler/internal/model"
	"gorm.io/gorm"
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
	var iv = &model.Intervals{
		IsCurrent:   1,
		IsCompleted: 0,
		Type:        t,
		TypeExtend:  ex,
	}
	interval, err := iv.First(i.db)
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
func (i *intervalsServant) CancelIntervalsCurrent(t int64, ex int64) error {
	var iv = &model.Intervals{
		IsCurrent:  0,
		Type:       t,
		TypeExtend: ex,
	}
	return iv.Updates(i.db, "IsCurrent")
}
func (i *intervalsServant) PickNewIntervals(t int64, ex int64) (*model.Intervals, error) {

	var iv = &model.Intervals{
		IsCurrent:   0,
		IsCompleted: 0,
		Type:        t,
		TypeExtend:  ex,
	}
	interval, err := iv.First(i.db)
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
