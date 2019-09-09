package designer

import (
	"canglong/speech_designer/common/interface/storage"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type DesignerServer struct {
	storage.Storage
	ip 		string
	port 	int
}

func (s *DesignerServer) Listen() {
	url := fmt.Sprintf("%s:%d",s.ip, s.port)

	logs.Debug("Robot Server Listen : http://%s/query\n", url)

	http.Handle("/query", s.newDesignerHandler())

	http.ListenAndServe(url,nil)
}

func New(s storage.Storage,ip string, port int) *DesignerServer {
	return &DesignerServer{s,ip, port}
}

