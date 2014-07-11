package kmdb

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"strconv"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func ProtobufNettyDecode(buf []byte, msg proto.Message) error {
	size, err := binary.ReadUvarint(bytes.NewBuffer(buf[:1]))
	if err != nil {
		return err
	}

	if size < 0 || int(size) >= len(buf) {
		return Error("Read message size error, size " + strconv.Itoa(int(size)))
	}

	return proto.Unmarshal(buf[1:size+1], msg)
}

func ProtobufNettyEncode(msg proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(msg)
	var rst []byte = make([]byte, 4, 4)
	err = binary.Write(rst, binary.LittleEndian, len(buf))
	rst = rst[3:4]
	append(rst, buf...)
	return rst, err
}
