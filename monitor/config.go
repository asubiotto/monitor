package monitor

import (
	"sync"
	"time"
)

type Config struct {
	// How often top section hits should be reported.
	reportInterval time.Duration
	// How many top sections should be shown.
	reportLimit int
	// If total traffic averaged over the number of sections goes over
	// trafficThreshold, report a traffic spike.
	trafficWindow    time.Duration
	trafficThreshold int
}

var (
	config    *Config
	configMtx sync.RWMutex

	// Our default config will default to the const values.
	defaultConfig = Config{
		reportInterval:   reportInterval,
		reportLimit:      reportLimit,
		trafficWindow:    trafficWindow,
		trafficThreshold: trafficThreshold,
	}
)

func SetConfig(cfg Config) {
	configMtx.Lock()
	defer configMtx.Unlock()
	config = &cfg
}

func GetConfig() *Config {
	// Lock just in case we have to set the config.
	configMtx.Lock()
	defer configMtx.Unlock()

	if config == nil {
		config = &defaultConfig
	}

	return config
}
