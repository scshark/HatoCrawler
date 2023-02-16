/**
 * @Author: scshark
 * @Description:
 * @File:  twitter
 * @Date: 1/21/23 3:35 PM
 */
package jinzhu

import (
	"HatoCrawler/internal/core"
	"HatoCrawler/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/big"
)
var (
	_ core.TweetService = (*tweetServant)(nil)
	_ core.TweetUserService = (*tweetUserServant)(nil)

)

type tweetServant struct {
	db *gorm.DB
}

type tweetUserServant struct {
	db *gorm.DB
}

func newTweetService(db *gorm.DB) core.TweetService {

	return &tweetServant{
		db: db,
	}
}
func newTweetUserService(db *gorm.DB) core.TweetUserService {

	return &tweetUserServant{
		db: db,
	}
}

func (t *tweetServant)CreateTweet(te *model.Twitter,data []model.Twitter) error{
	_, err := te.Create(t.db, data)
	return err
}
func (t *tweetServant)TweetIsExists(id *big.Int) bool{

	// tweet := &model.Twitter{
	// 	BModel : &model.BModel{
	// 		ID:id,
	// 	},
	// }
	tweet := &model.Twitter{
		Model : &model.Model{
			ID:id.Int64(),
		},
	}
	err := tweet.Exists(t.db)
	if err != nil{
		return  false
	}else {
		return  true
	}
}

func (t *tweetServant)GetTweetByUserId(userId int64,c *model.ConditionsT)(*model.Twitter,error)  {
	tweet := &model.Twitter{
		HtTwitterUserId:userId,
	}
	return  tweet.First(t.db,c)

}
func (t *tweetServant)CountTweetByUserId(userId int64,c *model.ConditionsT)(int64,error)  {
	tweet := &model.Twitter{
		HtTwitterUserId:userId,
	}
	return  tweet.Count(t.db,c)

}





func (t *tweetUserServant)CreateTweetUser(tu *model.TwitterUser) (*model.TwitterUser, error) {
	return tu.Create(t.db)
}

func (t *tweetUserServant)UpdateTweetUser(tu *model.TwitterUser) error {
	return tu.Update(t.db)
}

func (t *tweetUserServant)GetTweetUserByTweetId(id string) *model.TwitterUser  {

	var tUser = &model.TwitterUser{
		TweetUserId:id,
	}
	user,err := tUser.First(t.db)

	if err != nil  {
		logrus.Errorf(" tweet user first err %s",err)
		return nil
	}

	return user
}

func (t *tweetUserServant)GetTweetUserForLoad(twUser *model.TwitterUser) *model.TwitterUser  {

	user,err := twUser.First(t.db)

	if err != nil  {
		logrus.Errorf(" tweet user first err %s",err)
		return nil
	}

	return user
}

