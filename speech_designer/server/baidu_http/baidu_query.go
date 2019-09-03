package baidu_http

import (
	"canglong/speech_designer/common/message"
	"canglong/speech_designer/common/status"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

type baiduAsrHandler struct {
	*BaiduAsrServer
}

func (h *baiduAsrHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request message.RequestBody
	var text string
	response := message.ReponseBody{Code:status.StatusOK, Reason:status.StatusText(status.StatusOK)}

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

		if request.NodeParent != 0 {
			asr, stat := BaiduAsrBodyDecode(request)
			if stat != status.StatusOK {
				logs.Error("Failed to Request Unmarshal xml :%s\n",status.StatusText(stat))
				return
			}
			text = asr.Interpretation.Input
		} else {
			text = ""
		}

		nd ,stat := h.QueryNodeByRequest(request.SpeechId, request.SpeechNumber, request.NodeParent, text)
		if stat != status.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			logs.Error("Failed to query1 : %s", status.StatusText(stat))
			response.Code = status.StatusInternalServerError
			response.Reason = status.StatusText(status.StatusInternalServerError)
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

func (s *BaiduAsrServer) newBaiduAsrHandler() http.Handler {
	return &baiduAsrHandler{s}
}