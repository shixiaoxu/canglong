package config

import (
	c "canglong/com/config"
)

type LogsParams struct {
	FileName 	string 		`json:"filename"`
	Level 		int			`json:"level"`
	MaxLines 	int			`json:"maxlines"`
	MaxSize 	int			`json:"maxsize"`
	Daily		bool		`json:"daily"`
	MaxDays		int			`json:"maxdays"`
	Color 		bool		`json:"color"`
}

type LogsConf struct {
	Adapter	string
	Params 	LogsParams
}

var(
	Logger = LogsConf{
		Adapter:"console",
		Params:LogsParams{
			FileName:"log/esl_server.log",
			Level:7,
			MaxLines:1000000,
			MaxSize:4<<28,   //1<<28 = 256M   4<<28 = 10M
			Daily:true,
			MaxDays:0,
			Color:true,
	    },
	}
)

func (l *LogsConf)loadConf(configer c.Configer) {
	strVal := configer.String("logs::adapter")
	if strVal != "" {
		Logger.Adapter = strVal
	}

	strVal = configer.String("logs::filename")
	if strVal != "" {
		Logger.Params.FileName = strVal
	}

	intVal, err := configer.Int("logs::level")
	if err == nil {
		Logger.Params.Level = intVal
	}

	intVal, err = configer.Int("logs::maxlines")
	if err == nil {
		Logger.Params.MaxLines = intVal
	}

	intVal, err = configer.Int("logs::maxsize")
	if err == nil {
		Logger.Params.MaxSize = intVal * 1024 * 1024
	}

	intVal, err = configer.Int("logs::maxdays")
	if err == nil {
		Logger.Params.MaxDays = intVal
	}

	boolVal, err := configer.Bool("logs::daily")
	if err != nil {
		Logger.Params.Daily = boolVal
	}

	boolVal, err = configer.Bool("logs::color")
	if err != nil {
		Logger.Params.Color = boolVal
	}
}


