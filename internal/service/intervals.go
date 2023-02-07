/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 1/6/23 3:02 PM
 */
package service

import (
	"SecCrawler/internal/model"
	"github.com/sirupsen/logrus"
)

func InitLivesIntervals(nextCur int64, LivesType int64, typeExtend int64) error {


	// 获取一个当前区间
	cur, err := ds.GetCurrentIntervals(LivesType,typeExtend)
	if err != nil {
		logrus.Errorf("获取当前区间失败 error %s",err)
		return err
	}
	if cur.Model == nil {
		iv := &model.Intervals{
			Begin:     nextCur,
			Over:      0,
			Type:      LivesType,
			IsCurrent: 1,
			TypeExtend: typeExtend,
		}
		_, err = ds.CreateIntervals(iv)
		return err
	}
	// 无需初始化 x>= y
	if cur.Begin >= nextCur {
		return nil
	}
	err = ds.CancelIntervalsCurrent(LivesType,typeExtend)
	if err != nil {
		logrus.Errorf("取消所有区间状态失败 error %s",err)
		return err
	}

	iv := &model.Intervals{
		Begin:     nextCur,
		Over:      cur.Begin,
		Type:      LivesType,
		IsCurrent: 1,
		TypeExtend: typeExtend,
	}
	_, err = ds.CreateIntervals(iv)
	return err
}

func GetLivesCursor(LivesType int64,typeExtend int64) int64{

	iv, err := ds.GetCurrentIntervals(LivesType,typeExtend)
	if err != nil {
		logrus.Errorf("获取当前区间失败 error %s",err)
	}
	return iv.Begin
}

func UpdateLivesCursor(LivesType int64,next int64,typeExtend int64) error {
	// 最后一个id 更新到区间
	iv, err := ds.GetCurrentIntervals(LivesType,typeExtend)

	if err != nil {
		logrus.Errorf("获取当前区间失败 error %s",err)
		return err
	}
	if iv.Begin > next {
		iv.Begin = next
		err = ds.UpdateIntervals(iv)
		if err != nil {
			logrus.Errorf("更新当前区间失败 error %s",err)
		}
	}
	return err

}