/**
 * @Author: scshark
 * @Description:
 * @File:  internal
 * @Date: 12/8/22 6:19 PM
 */
package internal

import (
	"SecCrawler/internal/conf"
	"SecCrawler/internal/service"
)

func Internal() {
	// init service
	service.Initialize()

	conf.SetupRedisEngine()

}