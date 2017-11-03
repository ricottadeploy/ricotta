package comms

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/vmihailenco/msgpack"
)

type Conn struct {
	conn net.Conn
}

func NewConn(conn net.Conn) Conn {
	return Conn{conn}
}

func (conn *Conn) Write(data []byte) {
	c := conn.conn
	lenBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenBytes, uint16(len(data)))
	c.Write(lenBytes)
	c.Write(data)
}

func (conn *Conn) Read() ([]byte, error) {
	lenBytes := make([]byte, 2)
	count, err := conn.conn.Read(lenBytes)
	if err != nil {
		return nil, err
	}
	if count != 2 {
		return nil, errors.New("Unexpected length")
	}
	length := binary.LittleEndian.Uint16(lenBytes)

	data := make([]byte, length)
	count1, err1 := conn.conn.Read(data)
	if err1 != nil {
		return nil, err
	}
	if count1 != int(length) {
		return nil, errors.New("Unexpected message size")
	}
	return data, nil
}

func (conn *Conn) SendMsgFrame(frame MsgFrame) error {
	bytes, err := msgpack.Marshal(frame)
	if err != nil {
		return errors.New("Error while marshalling message")
	}
	conn.Write(bytes)
	return nil
}

func (conn *Conn) ReadMsgFrame() (MsgFrame, error) {
	frame := MsgFrame{}
	bytes, err := conn.Read()
	if err != nil {
		return frame, errors.New("Error while reading message")
	}
	err1 := msgpack.Unmarshal(bytes, &frame)
	if err1 != nil {
		return frame, errors.New("Error while unmarshalling MsgFrame")
	}
	return frame, nil
}
