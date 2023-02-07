/**
 * @Author: scshark
 * @Description:
 * @File:  intervals
 * @Date: 12/30/22 1:31 PM
 */
package model

import (
	"gorm.io/gorm"
)

type Intervals struct {
	*Model
	Begin       int64 `json:"begin"`
	Over        int64 `json:"over"`
	Type        int64    `json:"type"`
	IsCompleted int64    `json:"is_completed"`
	IsCurrent   int64    `json:"is_current"`
	TypeExtend  int64    `json:"type_extend"`
}

func (i *Intervals) Create(db *gorm.DB) (*Intervals, error) {
	err := db.Create(&i).Error
	return i, err
}

func (i *Intervals) First(db *gorm.DB) (*Intervals, error) {
	var iv Intervals
	if i.Model != nil && i.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", i.Model.ID, 0)
	} else if i.IsCurrent > 0 || i.Type > 0 {
		db = db.Where("type= ? AND is_current= ? AND is_completed= ?", i.Type, i.IsCurrent, i.IsCompleted).Order("`over` desc")
	}


	err := db.Limit(1).Find(&iv).Error
	return &iv, err
}
func (i *Intervals) Update(db *gorm.DB) error {
	return db.Model(&Intervals{}).Where("id = ? AND is_del = ?", i.Model.ID, 0).Save(i).Error
}
func (i *Intervals) Updates(db *gorm.DB, s ...string) error {

	db = db.Model(&Intervals{})
	if i.Type > 0 {
		db = db.Where("type= ? ", i.Type)
	}
	return db.Select(s).Updates(i).Error
}
