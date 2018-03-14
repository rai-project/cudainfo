package cudainfo

import (
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.WithField("pkg", "cudainfo")
		err := LoadUVM()
		if err != nil {
			log.WithError(err).Error("Failed to load uvm")
		}
		initNVMLLibrary()
		cnt, err := GetDeviceCount()
		if err != nil {
			log.WithError(err).Error("Was not able to query devices")
		}
		log.WithField("device_count", cnt).
			Infof("%d devices where found on the system", cnt)
	})
}
