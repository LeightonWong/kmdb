package kmdb

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
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
}

type StatusReply string

type IntegerReply int

type ErrorReply string

type BulkReply []byte

type MultiBulkReply struct {
	size    int
	replies []BulkReply
}
