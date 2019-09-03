package message

import "canglong/speech_designer/common/interface/storage"

type RequestBody struct {
	SpeechId 			int64 		`json:"SpeechId"`
	SpeechNumber		string		`json:"SpeechNumber"`
	NodeParent			int64		`json:"NodeParent"`
	NodeBody	 		string		`json:"NodeBody"`
}

type ReponseBody struct {
	Code 		int
	Reason 		string
	Node 		storage.S_Node
}
