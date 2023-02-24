/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 12/30/22 1:54 PM
 */
package core

import "HatoCrawler/internal/model"

type IntervalsService interface {

	//  创建区间
	CreateIntervals(i *model.Intervals) (*model.Intervals, error)
	//  取消所有的区间激活状态
	CancelIntervalsCurrent(t int64,ex int64) error
	//  获取当前激活区间
	GetCurrentIntervals(t int64,ex int64) (*model.Intervals, error)
	//  修改区间
	UpdateIntervals(i *model.Intervals) error

	GetPrevIntervals(t int64,ex int64) (*model.Intervals, error)

}