package kmdb

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"log"
	"strconv"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func ProtobufDecode(buf []byte, msg proto.Message) error {
	size, err := binary.ReadUvarint(bytes.NewBuffer(buf[:1]))
	if err != nil {
		return err
	}

	if size < 0 || int(size) >= len(buf) {
		return Error("Read message size error, size " + strconv.Itoa(int(size)))
	}

	return proto.Unmarshal(buf[1:size+1], msg)
}

func ProtobufEncode(msg proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return buf, err
	}
	var sizeByte []byte
	sizeBuf := bytes.NewBuffer(sizeByte)
	err = binary.Write(sizeBuf, binary.LittleEndian, int32(len(buf)))
	if err != nil {
		log.Println("Write data size error", err)
		return sizeByte, err
	}
	result := sizeBuf.Bytes()[:1]
	result = append(result, buf...)
	//log.Println("littleendian", result)
	return result, err
}
