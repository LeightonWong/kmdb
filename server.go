package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wlsailor/kmdb"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var config *Config
var kmdb *KMDB

func main() {
	runtime.GOMAXPROCS(8)
	welcome()
	init()
}

func init() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}
	/*
		if os.Args[0] == '-d' {
			isDaemon = true
		}
	*/
	config = kmdb.Config.Init(os.Args[1])
	kmdb = KMDB.Open(config)
	if kmdb == nil {
		log.Fatalf("Could not open kmdb.")
		os.Exit(1)
	}
	CheckPidFile()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", config.Listen.Ip+":"+config.Listen.Port)
	listener, _ := net.Listen("tcp", tcpAddr)
	log.Printf("KMDB start to accepting connections.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleClient(conn)
	}
	WritePid()
	log.Printf("KMDB server started.\n")
}

func HandleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 255)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("kmdb protocal error", err)
			return
		}
		comm := kmdb.Command{}
		err = json.Unmarshal(buf, &comm)
		if err != nil {
			log.Println("received msg parse error", buf)
		}
		if comm.Sync == false {
			go HandleOperation(conn, &comm)
		} else {
			HandleOperation(conn, &comm)
		}
	}
}

func HandleOperation(conn net.Conn, comm *kmdb.Command) {
	var err error
	switch comm.Type {
	case "Get":
		value, err := kmdb.Get(comm.Key)
		HandleError(conn, err)
		result := kmdb.Result{0, "Success", value}
		bytes, err := json.Marshal(result)
		if err != nil {
			log.Println("Serlize get result error, result:", result, " error:", err)
		}
		_, err = conn.Write(bytes)

	case "Put":
		err = kmdb.Put(comm.Key, comm.Value)
		HandleError(conn, err)
		result := json.Marshal(kmdb.Result{0, "Success", []byte("1")})
		_, err = conn.Write(result)
	case "Del":
		err = kmdb.Del(comm.Key)
		HandleError(conn, err)
	default:
		fmt.Printf("Unsupported command type:%s", comm.Type)
	}
}

func HandleError(conn net.Conn, err error) {
	if err != nil {
		log.Println("Conn ", conn.RemoteAddr(), " error:", error)
		conn.Close()
	}
}

func welcome() {
	fmt.Printf("kmdb %s\n", kmdb.KMDB_VERSION)
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
	size, err := file.Read(pid)
	if err != nil {
		log.Panic(err)
	}
	return strconv.Atoi(string(pid))
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
