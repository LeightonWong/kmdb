package kmdb

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Listen struct {
		Ip   string
		Port int
	}

	Store struct {
		Dir     string
		PidFile string
	}

	Type struct {
		Primary bool
		SlaveOf string
	}
}

func LoadConfig(file string) *Config {
	goconfig.LineBreak = "\n"
	goconfig.PrettyFormat = true
	configFile, err := goconfig.LoadConfigFile(file)
	if err != nil {
		log.Fatalf("Read config %s error, exit.", file)
		os.Exit(0)
	}
	var config Config = Config{}
	config.Listen.Ip = configFile.MustValue("listen", "ip", "0.0.0.0")
	config.Listen.Port = configFile.MustInt("listen", "port", 2222)
	config.Store.Dir = configFile.MustValue("store", "dir", "/home/wanglei/kmdb")
	config.Store.PidFile = configFile.MustValue("store", "pidfile", "/tmp/kmdb-"+strconv.Itoa(config.Listen.Port)+".pid")
	config.Type.Primary = configFile.MustBool("type", "primary", true)
	if !config.Type.Primary {
		config.Type.SlaveOf = configFile.MustValue("type", "slaveOf")
	}
	return &config
}
