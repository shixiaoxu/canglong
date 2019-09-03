package storage

import (
	"canglong/speech_designer/common/status"
)

type Storage interface {
	Init() (status status.Status)
	DeInit() (status status.Status)

	QueryNodeByRequest(speech int64, number string, parent int64, text string) (en S_Node, stat status.Status)
}
