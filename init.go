package cudainfo

import (
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log *logrus.Entry
)

func init() {
	config.OnInit(func() {
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

func initNVML() error {
	if err := os.Setenv("CUDA_DISABLE_UNIFIED_MEMORY", "1"); err != nil {
		return err
	}
	if err := os.Setenv("CUDA_CACHE_DISABLE", "1"); err != nil {
		return err
	}
	if err := os.Unsetenv("CUDA_VISIBLE_DEVICES"); err != nil {
		return err
	}
	initNVMLLibrary()
	return nil
}
