package initial

import "canglong/esl_server/config"

func initConfig() {
	config.LoadConfig("ini", "conf/esl_server.conf")
}
