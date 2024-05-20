package notify

import "github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"

type Message string

type Notification[T NotificationPayload] struct {
	Msg     Message       `json:"msg"`
	Meta    *rpc.MetaData `json:"meta"`
	Payload T             `json:"payload"`
}

type NotificationPayload interface {
	[]NotificationPayloadAppearance |
		[]NotificationPayloadChunkWritten |
		NotificationPayloadChunkWritten |
		string
}
