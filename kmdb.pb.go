// Code generated by protoc-gen-go.
// source: kmdb.proto
// DO NOT EDIT!

/*
Package kmdb is a generated protocol buffer package.

It is generated from these files:
	kmdb.proto

It has these top-level messages:
	Command
	Result
*/
package kmdb

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type CommandType int32

const (
	CommandType_PUT CommandType = 0
	CommandType_GET CommandType = 1
	CommandType_DEL CommandType = 2
)

var CommandType_name = map[int32]string{
	0: "PUT",
	1: "GET",
	2: "DEL",
}
var CommandType_value = map[string]int32{
	"PUT": 0,
	"GET": 1,
	"DEL": 2,
}

func (x CommandType) Enum() *CommandType {
	p := new(CommandType)
	*p = x
	return p
}
func (x CommandType) String() string {
	return proto.EnumName(CommandType_name, int32(x))
}
func (x *CommandType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CommandType_value, data, "CommandType")
	if err != nil {
		return err
	}
	*x = CommandType(value)
	return nil
}

type Command struct {
	Type             *CommandType `protobuf:"varint,1,req,name=type,enum=kmdb.CommandType" json:"type,omitempty"`
	Key              []byte       `protobuf:"bytes,2,req,name=key" json:"key,omitempty"`
	Value            []byte       `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
	Sync             *bool        `protobuf:"varint,4,req,name=sync" json:"sync,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}

func (m *Command) GetType() CommandType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CommandType_PUT
}

func (m *Command) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Command) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Command) GetSync() bool {
	if m != nil && m.Sync != nil {
		return *m.Sync
	}
	return false
}

type Result struct {
	Code             *int32  `protobuf:"varint,1,req,name=code" json:"code,omitempty"`
	Msg              *string `protobuf:"bytes,2,req,name=msg" json:"msg,omitempty"`
	Rst              []byte  `protobuf:"bytes,3,opt,name=rst" json:"rst,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}

func (m *Result) GetCode() int32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *Result) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *Result) GetRst() []byte {
	if m != nil {
		return m.Rst
	}
	return nil
}

func init() {
	proto.RegisterEnum("kmdb.CommandType", CommandType_name, CommandType_value)
}