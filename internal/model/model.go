/**
 * @Author: scshark
 * @Description:
 * @File:  model
 * @Date: 12/12/22 4:42 PM
 */
package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)
type ConditionsT map[string]interface{}

type Model struct {
	ID int64 `gorm:"primary_key" json:"id"`
	CreatedOn  int64                 `json:"created_on"`
	ModifiedOn int64                 `json:"modified_on"`
	DeletedOn  int64                 `json:"deleted_on"`
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"`
}


func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	nowTime := time.Now().Unix()

	tx.Statement.SetColumn("created_on", nowTime)
	tx.Statement.SetColumn("modified_on", nowTime)
	return
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	if !tx.Statement.Changed("modified_on") {
		tx.Statement.SetColumn("modified_on", time.Now().Unix())
	}

	return
}
