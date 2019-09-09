package designer

import (
	"canglong/speech_designer/common/interface/storage"
	"canglong/speech_designer/common/status"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

type designerHandler struct {
	*DesignerServer
}

func (h *designerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request storage.RequestBody
	response := storage.ReponseBody{Code:status.StatusOK, Reason:status.StatusText(status.StatusOK)}

	if r.Method == http.MethodPost {
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logs.Error("Failed to ReadAll : %s\n", err.Error())
			return
		}

		logs.Debug(string(req))

		err = json.Unmarshal(req, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logs.Error("Failed to Request Unmarshal Json : %s\n", err.Error())
			return
		}

		nd ,stat := h.QueryNodeByRequest(request)
		if stat != status.StatusOK {
			logs.Error("Failed to query1 : %s", status.StatusText(stat))
			response.Code = int(stat)
			response.Reason = status.StatusText(stat)
		} else {
			response.Node = nd
		}

		json, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logs.Error("Failed to Response Marshal : %s", err.Error())
		}

		logs.Debug(string(json))

		w.Write([]byte(json))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *DesignerServer) newDesignerHandler() http.Handler {
	return &designerHandler{s}
}