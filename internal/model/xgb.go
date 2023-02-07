/**
 * @Author: scshark
 * @Description:
 * @File:  xgb
 * @Date: 1/10/23 12:28 PM
 */
package model

import "gorm.io/gorm"

type Xgb struct {
	*Model
	Title string `json:"title"`
	Summary string `json:"summary"`
	Image string `json:"image"`
	LiveCreatedAt int64 `json:"live_created_at"`
	SubjIds string `json:"subj_ids"`
	Uri string `json:"uri"`
	Tags string `json:"tags"`
	OriginaUrl string `json:"origina_url"`
	Source string `json:"source"`
}
func (x *Xgb) Create (db *gorm.DB,items []Xgb) (*Xgb,error){
	err := db.Model(&x).Create(items).Error
	return  x,err
}

func (x *Xgb) First(db *gorm.DB) (*Xgb, error) {
	var dh Xgb
	if x.Model != nil && x.Model.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", x.Model.ID, 0)
	}
	err := db.Limit(1).Find(&dh).Error
	return &dh, err
}

func (x *Xgb) Update(db *gorm.DB) error{
	return db.Model(&Intervals{}).Where("id = ? AND is_del = ?", x.Model.ID, 0).Omit("id").Save(x).Error
}