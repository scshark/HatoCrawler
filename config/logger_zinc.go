/**
 * @Author: scshark
 * @Description:
 * @File:  logger_zinc
 * @Date: 2/2/23 1:13 PM
 */
package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"time"
)

type zincLogData struct {
	Time    time.Time     `json:"time"`
	Level   logrus.Level  `json:"level"`
	Message string        `json:"message"`
	Data    logrus.Fields `json:"data"`
}

type zincLogIndex struct {
	Index map[string]string `json:"index"`
}

type zincLogHook struct {
	host     string
	index    string
	user     string
	password string
}

func (h *zincLogHook) Fire(entry *logrus.Entry) error {
	index := &zincLogIndex{
		Index: map[string]string{
			"_index": h.index,
		},
	}
	indexBytes, _ := json.Marshal(index)

	data := &zincLogData{
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
		Data:    entry.Data,
	}
	dataBytes, _ := json.Marshal(data)

	logStr := string(indexBytes) + "\n" + string(dataBytes) + "\n"
	client := resty.New()

	if _, err := client.SetDisableWarn(true).R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(h.user, h.password).
		SetBody(logStr).
		Post(h.host); err != nil {

		fmt.Println(err.Error())
	}
	return nil
}

func (h *zincLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newZincLogHook() *zincLogHook {
	return &zincLogHook{
		host:     Cfg.LoggerZinc.Endpoint()+ "/es/_bulk",
		index:    Cfg.LoggerZinc.Index,
		user:     Cfg.LoggerZinc.User,
		password: Cfg.LoggerZinc.Password,
	}
}
