package main

import (
	//"encoding/json"
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"flag"
	"fmt"
	km "github.com/wlsailor/kmdb"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
)

var config *km.Config
var db *km.KMDB
var configFile = flag.String("conf", "kmdb.conf", "config file path")

func main() {
	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Ldate)
	welcome()
	Init()
	defer db.Close()
}

func Init() {
	/*
		if len(os.Args) < 2 {
			usage()
			os.Exit(0)
		}
	*/
	flag.Parse()
	config = km.LoadConfig(*configFile)
	db = km.Open(config)
	if db == nil {
		log.Fatalf("Could not open kmdb.")
		os.Exit(1)
	}
	CheckPidFile()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", config.Listen.Ip+":"+strconv.Itoa(config.Listen.Port))
	listener, _ := net.Listen(tcpAddr.Network(), tcpAddr.String())
	fmt.Printf("KMDB initial successful, listen on %s:%d\n", config.Listen.Ip, config.Listen.Port)
	fmt.Printf("	Database Dir %s\n", config.Store.Dir)
	fmt.Printf("	Process ID Saved at %s\n", config.Store.PidFile)
	fmt.Printf("	DB is primary : %v\n", config.Type.Primary)
	log.Printf("kmdb start to accepting connections.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleClient(conn)
	}
	WritePid()
	log.Printf("km server started.\n")
}

func HandleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 128, 255)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed from remote")
			} else {
				log.Printf("kmdb protocal error %s", err)
			}
			defer conn.Close()
			return
		}
		comm := km.Command{}
		//err = json.Unmarshal(buf, &comm)
		size, err := binary.ReadUvarint(bytes.NewBuffer(buf[:1]))
		log.Println("Read buf size", len(buf), "proto buf length", size)
		if err != nil {
			log.Println("Message size error", buf)
			continue
		}
		message := buf[1 : size+1]
		err = proto.Unmarshal(message, &comm)
		if err != nil {
			log.Println("Income message parse failed, Msg:", message)
			log.Println(err)
			buf = make([]byte, 128, 255)
			continue
		}
		buf = make([]byte, 128, 255)
		if *comm.Sync == false {
			go HandleOperation(conn, &comm)
		} else {
			HandleOperation(conn, &comm)
		}
	}
}

func HandleOperation(conn net.Conn, comm *km.Command) {
	var err error
	switch *comm.Type {
	case km.CommandType_GET:
		value, err := db.Get(comm.Key)
		result := HandleError(conn, err)
		if err == nil {
			result.Rst = *value
		}
		bytes, err := km.ProtobufNettyEncode(result)
		if err != nil {
			log.Println("Serlize get result error, result:", result, " error:", err)
		}
		_, err = conn.Write(bytes)

	case km.CommandType_PUT:
		err = db.Put(comm.Key, comm.Value)
		rst := HandleError(conn, err)
		result, err := km.ProtobufNettyEncode(rst)
		_, err = conn.Write(result)
	case km.CommandType_DEL:
		err = db.Del(comm.Key)
		rst := HandleError(conn, err)
		bytes, err := km.ProtobufNettyEncode(rst)
		_, err = conn.Write(bytes)
	default:
		result := &km.Result{Code: proto.Int(-1), Msg: proto.String("Unsupported command")}
		bytes, err := km.ProtobufNettyEncode(result)
		if err != nil {
			log.Println("Serlize get result error, result:", result, " error:", err)
		}
		_, err = conn.Write(bytes)
		HandleError(conn, err)
	}
}

func HandleError(conn net.Conn, err error) *km.Result {
	if err != nil {
		log.Println("Conn ", conn.RemoteAddr(), " error:", err)
		conn.Close()
		return &km.Result{Code: proto.Int(1), Msg: proto.String(err.Error())}
	} else {
		return &km.Result{COde: proto.Int(0), Msg: proto.String("Success")}
	}
}

func welcome() {
	fmt.Printf("kmdb %s\n", km.KMDB_VERSION)
	fmt.Printf("Copyright (C) 2014-2015 kongming-inc.com\n")
	fmt.Printf("\n")
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("	%s [-d] /path/to/kmdb.conf\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("	-d	run as daemon\n")
}

func WritePid() {
	pid := os.Getpid()
	pidFile := config.Store.PidFile
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_RDWR, 0755)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	file.WriteString(strconv.Itoa(pid))
}

func ReadPid() int {
	pidFile := config.Store.PidFile
	file, err := os.OpenFile(pidFile, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	var pid []byte = make([]byte, 10, 10)
	_, err = file.Read(pid)
	if err != nil {
		log.Panic(err)
	}
	r, err := strconv.Atoi(string(pid))
	if err != nil {
		return 0
	} else {
		return r
	}
}

func CheckPidFile() {
	pidFile := config.Store.PidFile
	if len(pidFile) > 0 {
		if _, err := os.Stat(pidFile); err == nil {
			fmt.Printf("Fatal error!Pidfile %s already exists!")
			os.Exit(1)
		}
	}
}

func RemovePidFile() {
	pidFile := config.Store.PidFile
	os.Remove(pidFile)
}
