package baidu_http

import (
	"canglong/speech_designer/common/interface/storage"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type BaiduAsrServer struct {
	storage.Storage
	ip 		string
	port 	int
}

func (s *BaiduAsrServer) Listen() {
	url := fmt.Sprintf("%s:%d",s.ip, s.port)

	logs.Debug("Robot Server Listen : http://%s/query\n", url)

	http.Handle("/query", s.newBaiduAsrHandler())

	http.ListenAndServe(url,nil)
}

func New(s storage.Storage,ip string, port int) *BaiduAsrServer {
	return &BaiduAsrServer{s,ip, port}
}

