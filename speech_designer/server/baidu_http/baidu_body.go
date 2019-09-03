package baidu_http

import (
	"canglong/speech_designer/common/message"
	"canglong/speech_designer/common/status"
	"encoding/base64"
	"encoding/xml"
	"github.com/astaxie/beego/logs"
)

type Input struct {
	XMLName 		xml.Name 	`xml:"input"`
	Mode 			string 		`xml:"mode,attr"`
}

type ExtJson struct{
	XMLName 		xml.Name	`xml:"extJson"`
	SnStartTime		string		`xml:"snStartTime"`
	SnStopTime		string		`xml:"snStopTime"`
	Speed			string		`xml:"speed"`
}

type Instance struct {
	XMLName 		xml.Name	`xml:"instance"`
	CallId			string		`xml:"callId"`
	LogId			string		`xml:"logId"`
	Rolecategory	string		`xml:"rolecategory"`
	CategotyId		string		`xml:"categotyId"`
	ExtJson			ExtJson		`xml:"extJson"`
}

type Interpretation struct {
	XMLName 	xml.Name	`xml:"interpretation"`
	Grammer 	string		`xml:"grammar,attr"`
	Confidence  string		`xml:"confidence,attr"`
	Instance 	Instance	`xml:"instance"`
	Input	 	string		`xml:"input"`
}

type Result struct {
	XMLName 			xml.Name 			`xml:"result"`
	Interpretation 		Interpretation 		`xml:"interpretation"`
}

func BaiduAsrBodyDecode(req message.RequestBody) (rc Result , stat status.Status) {
	stat = status.StatusOK

	data, err := base64.StdEncoding.DecodeString(req.NodeBody)
	if err != nil {
		logs.Error("Failed to base64 : %s", err)
		stat = status.StatusBadRequest
		return rc, stat
	}
	logs.Debug(string(data))
	err = xml.Unmarshal(data, &rc)
	if err != nil {
		stat = status.StatusBadRequest
		return rc, stat
	}

	return rc, stat
}