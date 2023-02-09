/**
 * @Author: scshark
 * @Description:
 * @File:  servic
 * @Date: 12/8/22 6:27 PM
 */
package service

import (
	"HatoCrawler/internal/core"
	"HatoCrawler/internal/dao"
)

var (
	ds core.DataService
)

func Initialize() {
	ds = dao.DataService()
}
