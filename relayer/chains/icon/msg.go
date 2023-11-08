package icon

import (
	"encoding/json"
)

const defaultStepLimit = 13610920010

type IconMessage struct {
	Params interface{}
	Method string
}

func (im *IconMessage) Type() string {
	return im.Method
}

func (im *IconMessage) MsgBytes() ([]byte, error) {
	return json.Marshal(im.Params)
}

func (icp *IconProvider) NewIconMessage(msg interface{}, method string) IconMessage {
	return IconMessage{
		Params: msg,
		Method: method,
	}
}
