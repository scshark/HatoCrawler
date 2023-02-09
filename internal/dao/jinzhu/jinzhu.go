/**
 * @Author: scshark
 * @Description:
 * @File:  jinzhu
 * @Date: 12/9/22 3:01 PM
 */
package jinzhu

import (
	"HatoCrawler/internal/conf"
	"HatoCrawler/internal/core"
)

var (
	_ core.DataService = (*dataServant)(nil)
)
type dataServant struct {
	core.JinseService
	core.WallStreetService
	core.IntervalsService
	core.DyhjwService
	core.XgbService
	core.TweetService
	core.TweetUserService
}
func NewDataService() core.DataService{

	// initialize CacheIndex if needed

	db := conf.MustGormDb()

	ds := &dataServant{
		JinseService:  newJinseService(db),
		WallStreetService:  newWallStreetService(db),
		IntervalsService:  newIntervalsService(db),
		DyhjwService:  newDyhjwService(db),
		XgbService:  newXgbService(db),
		TweetService:  newTweetService(db),
		TweetUserService:  newTweetUserService(db),
	}

	return ds
}