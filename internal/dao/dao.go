/**
 * @Author: scshark
 * @Description:
 * @File:  dao
 * @Date: 12/8/22 6:31 PM
 */
package dao

import (
	"SecCrawler/internal/core"
	"SecCrawler/internal/dao/jinzhu"
	"sync"
)

var (
	OnceDs sync.Once
	ds core.DataService
)
func DataService() core.DataService{


	OnceDs.Do(func() {
		ds = jinzhu.NewDataService()
	})

	return ds
}