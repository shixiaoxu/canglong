package mysql

import (
	"canglong/speech_designer/common/interface/storage"
	"canglong/speech_designer/common/status"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"time"
)

type storageMysql struct {
	driver 		string
	driverUrl 	string
	db 			*sqlx.DB
}

func (s *storageMysql) QueryNodeByRequest(speech int64, number string, parent int64, text string) (en storage.S_Node, stat status.Status) {
	var node []storage.S_Node
	stat = status.StatusNotFound
	en = storage.S_Node{
		SN_ID:parent, S_ID:1,
		SN_Parent:parent,
		SN_Action:"play_and_detect_speech.file",
		SN_Text:"对不起请说普通话或找管理设置话术",
		SN_File:"test_error.mp3",
		SN_Argc:0,
		SN_Argv:"",
	}

	defer func() {
		err := recover()
		if err != nil {
			logs.Error("Failed to mysql query : %s", err)
		}
	}()

	if parent >= 0 {
		tm := time.Now()
		date := fmt.Sprintf("%d-%02d-%02d",tm.Year(), tm.Month(),tm.Day())
		time := fmt.Sprintf("%02d:%02d:%02d", tm.Hour(), tm.Minute(), tm.Second())
		sql := fmt.Sprintf(SqlQueryNodeFormartByRequest, text, parent, number, date, time)
		logs.Debug(sql)
		err := s.db.Select(&node, sql)
		if err != nil {
			logs.Error(err.Error())
			return en, stat
		}

		stat = status.StatusOK
		en = node[0]
	} else {
		return en, stat
	}

	return en, stat
}

func (s *storageMysql)Init() (status.Status) {
	db ,err := sqlx.Open(s.driver, s.driverUrl)
	if err != nil {
		logs.Error(err)
		return status.StatusMysqlConnectFailed
	}

	s.db = db
	return status.StatusOK
}

func (s *storageMysql)DeInit() (status.Status) {
	s.db.Close()
	return status.StatusOK
}

func New(driver, driverUrl string) storage.Storage {
	return &storageMysql{driver, driverUrl, nil}
}