package config

import (
	c "canglong/com/config"
	"fmt"
)

func LoadConfig(adapter, file string) {
	configer ,err := c.NewConfig(adapter, file)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config : %s", err.Error()))
	}

	Logger.loadConf(configer)
	Esl.LoadConf(configer)
	return
}

