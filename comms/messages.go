package comms

import (
	"errors"

	"github.com/vmihailenco/msgpack"
)

type MsgType uint8

const (
	MsgTypeInfo MsgType = iota
	MsgTypeCommand
)

type MsgFrame struct {
	Type MsgType
	Data []byte
}

type InfoMessage struct {
	Text string
}

func NewInfoMessage(text string) (MsgFrame, error) {
	msg := InfoMessage{Text: text}
	msgBytes, err := msgpack.Marshal(msg)
	if err != nil {
		return MsgFrame{}, errors.New("Error while marshalling message")
	}
	frame := MsgFrame{Type: MsgTypeInfo, Data: msgBytes}
	return frame, nil
}

func ToInfoMessage(data []byte) (InfoMessage, error) {
	msg := InfoMessage{}
	err := msgpack.Unmarshal(data, &msg)
	if err != nil {
		return msg, errors.New("Error while unmarshalling InfoMessage")
	}
	return msg, nil
}
