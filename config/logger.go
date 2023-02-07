/**
 * @Author: scshark
 * @Description:
 * @File:  logger
 * @Date: 2/2/23 2:52 PM
 */
package config

import (
	"github.com/sirupsen/logrus"
	"io"
)


func SetupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(Cfg.Logger.logLevel())

	hook := newZincLogHook()
	logrus.SetOutput(io.Discard)
	logrus.AddHook(hook)
}
