package main

import (
	"canglong/speech_designer/common/status"
	sto "canglong/speech_designer/storage/mysql"
	hser "canglong/speech_designer/server/baidu_http"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	st := sto.New("mysql" ,"root:123qwe@tcp(192.168.85.132)/canglong")
	stat := st.Init()
	if stat != status.StatusOK {
		panic(status.StatusText(stat))
	}
	defer st.DeInit()

	logs.Debug("Robot Storage Connect success")

	hser.New(st,"192.168.43.94",8000).Listen()
}
