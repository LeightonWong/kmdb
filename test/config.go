package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
)

type Config struct {
	Listen struct {
		Ip   string
		Port int
	}
}

func main() {
	config, err := goconfig.LoadConfigFile("kmdb.cfg")
	if err != nil {
		panic(err)
	}
	fmt.Printf("KMDB start up listen on : %s", config.MustValue("listen", "ip"))
}
