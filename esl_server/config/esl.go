package config

import (
	c "canglong/com/config"
)

type EslConf struct {
	Host 		string
	Password 	string
	Port 		int
	Timeout 	int
}

var (
	Esl	= EslConf {
		Host:"127.0.0.1",
		Password:"123qwe",
		Port:8021,
		Timeout:3000,
	}
)

func (e *EslConf)LoadConf(configer c.Configer) {
	strVal := configer.String("freeswitch::host")
	if strVal != "" {
		Esl.Host = strVal
	}

	strVal = configer.String("freeswitch::password")
	if strVal != "" {
		Esl.Password = strVal
	}

	intVal, err:= configer.Int("freeswitch::port")
	if err == nil {
		Esl.Port = intVal
	}

	intVal, err = configer.Int("freeswitch::timeout")
	if err == nil {
		Esl.Timeout = intVal
	}
}

