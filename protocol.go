package kmdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

const (
	ASTERISK_BYTE byte = '*'
	DOLLAR_BYTE   byte = '$'
	PLUS_BYTE     byte = '+'
	MINUS_BYTE    byte = '-'
	COLON_BYTE    byte = ':'
	CR_BYTE       byte = '\r'
	LF_BYTE       byte = '\n'
)

type Operation int

const (
	PING Operation = iota
	GET
	PUT
	DEL
)

type Command struct {
	op   Operation
	args []byte
	sync bool
}

type CommandBatch struct {
	comms []Command
	sync  bool
}

type StatusReply string

type IntegerReply int

type ErrorReply string

type BulkReply []byte

type MultiBulkReply struct {
	size    int
	replies []BulkReply
}

func DecodeCommand(buf []byte, comm *Command) error {
	if len(buf) == 0 {
		return Error("Empty data")
	}
	if buf[0] != ASTERISK_BYTE {
		return Error("Wrong Command Format " + string(buf))
	}

	oplen, err := binary.ReadUvarint(bytes.NewBuffer(buf[1:1]))

}
