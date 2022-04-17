package helpers

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var log = logrus.New()
var doOnce sync.Once

func GetLogger() *logrus.Logger {
	doOnce.Do(func() {
		log.WithFields(map[string]interface{}{
			"app_name": "tasmota_backup",
			"version":  "debug",
		})
		log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
		log.SetReportCaller(true)
	})
	return log
}
