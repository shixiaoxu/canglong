package storage

import (
	"canglong/speech_designer/common/status"
)

type Storage interface {
	Init() (status status.Status)
	DeInit() (status status.Status)

	QueryNodeByRequest(body RequestBody) (en S_Node, stat status.Status)
}
