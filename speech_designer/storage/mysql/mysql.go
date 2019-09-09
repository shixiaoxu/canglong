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


func (s *storageMysql) QueryNodeByRequest(body storage.RequestBody) (en storage.S_Node, stat status.Status) {
	var sql string
	var dbnode []storage.S_Node
	tm := time.Now()
	date := fmt.Sprintf("%d-%02d-%02d",tm.Year(), tm.Month(),tm.Day())
	time := fmt.Sprintf("%02d:%02d:%02d", tm.Hour(), tm.Minute(), tm.Second())

	defer func() {
		err := recover()
		if err != nil {
			logs.Error("Failed to mysql query : %s", err)
			en = storage.S_Node{}
			stat = status.StatusDBQueryRowNotFound
		}
	}()

	if body.QueryMethod == "Next" {
		sql = fmt.Sprintf(SqlQueryNextFormat, body.NodeTrigger, body.NodeParent, body.SpeechNumber, date, time)
	} else if body.QueryMethod == "Frist" {
		sql = fmt.Sprintf(SqlQueryFristFormat, body.SpeechNumber, date, time)
	} else if body.QueryMethod == "Parent" {
		sql = fmt.Sprintf(SqlQueryParentFormat, body.SpeechNumber, body.NodeParent, date, time)
	} else {
		en = storage.S_Node{}
		stat = status.StatusDBQueryMethod
		return en, stat
	}

	logs.Debug(sql)
	err := s.db.Select(&dbnode, sql)
	if err != nil {
		en = storage.S_Node{}
		stat = status.StatusDBQueryRowNotFound
	} else {
		en = dbnode[0]
		stat = status.StatusOK
	}

	return en, stat
}

func (s *storageMysql)Init() (status.Status) {
	db ,err := sqlx.Open(s.driver, s.driverUrl)
	if err != nil {
		logs.Error(err)
		return status.StatusDBConnectFailed
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