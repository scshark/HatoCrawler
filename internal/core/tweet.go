/**
 * @Author: scshark
 * @Description:
 * @File:  tweet
 * @Date: 1/17/23 4:16 PM
 */
package core

import (
	"HatoCrawler/internal/model"
	"math/big"
)

type TweetService interface {
	CreateTweet(te *model.Twitter,data []model.Twitter) error
	TweetIsExists(id *big.Int) bool
	GetTweetByUserId(userId int64,c *model.ConditionsT) (*model.Twitter,error)
}

type TweetUserService interface {
	CreateTweetUser(tu *model.TwitterUser) (*model.TwitterUser, error)
	UpdateTweetUser(tu *model.TwitterUser) error
	GetTweetUserByTweetId(id string) *model.TwitterUser
	GetTweetUserForLoad(twUser *model.TwitterUser) *model.TwitterUser
}