/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 1/6/23 3:02 PM
 */
package service

import (
	"HatoCrawler/internal/model"
	"github.com/sirupsen/logrus"
)

func InitLivesIntervals(nextCur int64, LivesType int64, typeExtend int64) error {


	overCursor,e := ds.GetIntervalsOverCursor(LivesType,typeExtend)

	if e !=nil {
		logrus.Errorf("初始化获取区间节点末位标记失败 GetIntervalsOverCursor error %s",e)
		return e
	}
	if overCursor == 0 {
		logrus.Infof("区间节点 末位 为 0 ，LivesType %v ",LivesType)
	}

	// 无需初始化 x>= y
	if overCursor >= nextCur {
		return nil
	}
	err := ds.CancelIntervalsCurrent(LivesType,typeExtend)
	if err != nil {
		logrus.Errorf("取消所有区间状态失败 error %s",err)
		return err
	}


	iv := &model.Intervals{
		Begin:     nextCur,
		Over:      overCursor,
		Type:      LivesType,
		IsCurrent: 1,
		TypeExtend: typeExtend,
	}
	_, err = ds.CreateIntervals(iv)
	if err != nil {
		logrus.Errorf("创建初始化区间失败 error %s",err)
		return err
	}

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