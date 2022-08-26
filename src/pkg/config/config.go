package config

import (
	"encoding/json"
	"litmus-init/types"
	"os"

	log "github.com/sirupsen/logrus"
)

func Load() types.Config {
	var config types.Config
	configData := os.Getenv("CONFIG_INIT_AS_CODE")

	err := json.Unmarshal([]byte(configData), &config)
	if err != nil {
		log.Fatalf("Fatal error: cannot load configuration as code\n%s\n", err)
	}

	return config
}
