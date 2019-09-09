package storage

type Speech struct {
	S_ID			int64
	S_BDate			string
	S_EDate			string
	S_BTime			string
	S_ETime			string
	S_Priority 		int64
	S_Node 			int64
	S_Number		string
	S_Enabled 		int64
	S_Description 	string
}

type S_Node struct {
	SN_ID 			int64  		`db:"sn_id"`
	SN_Parent		int64		`db:"sn_parent"`
	SN_Action		string		`db:"sn_action"`
	SN_File			string		`db:"sn_file"`
	SN_Text			string		`db:"sn_text"`
	SN_Argc			int64		`db:"sn_argc"`
	SN_Argv			string		`db:"sn_argv"`
	SN_Description	string		`db:"sn_description"`
	S_ID			int64		`db:"s_id"`
	S_Priority		int64		`db:"s_priority"`
	SNT_Trigger 	string		`db:"snt_trigger"`
}

type SN_Trigger struct {
	SNT_ID			int64
	SNT_Trigger		string
	SN_ID			int64
}

type RequestBody struct {
	SpeechId 			int64 		`json:"SpeechId"`
	SpeechNumber		string		`json:"SpeechNumber"`
	NodeParent			int64		`json:"NodeParent"`
	NodeTrigger	 		string		`json:"NodeTrigger"`
	QueryMethod			string 		`json:"QueryMethod"`
}

type ReponseBody struct {
	Code 		int
	Reason 		string
	Node 		S_Node
}