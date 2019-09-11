package initial

import (
	"canglong/com/logs"
)

func initLog() {
	logs.SetLogger(logs.AdapterConsole, `{"level":7,"color":true}`)
	//logs.SetLogger(logs.AdapterFile, `{"filename":"esl_server.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
}