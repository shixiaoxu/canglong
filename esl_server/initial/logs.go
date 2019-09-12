package initial

import (
	"canglong/esl_server/config"
	"canglong/com/logs"
	"encoding/json"
	"fmt"
)

func initLog() {
	jsonP, err := json.Marshal(config.Logger.Params)
	if err != nil {
		panic(fmt.Sprintf("Failed to init log : %s", err.Error()))
	}

	if config.Logger.Adapter == "file" {
		logs.SetLogger(logs.AdapterFile, string(jsonP))
	} else {
		logs.SetLogger(logs.AdapterConsole, string(jsonP))
	}

	jsonL, err := json.Marshal(config.Logger)
	logs.Debug(string(jsonL))
}
