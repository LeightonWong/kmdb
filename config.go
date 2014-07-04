package kmdb

import (
	"github.com/Unknwon/goconfig"
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

func Init(file string) *Config {
	configFile, err := goconfig.LoadConfigFile(file)
	var config Config = Config{}
	config.Listen.Ip = configFile.MustValue("listen", "ip", "0.0.0.0")
	config.Listen.Port = configFile.MustValue("listen", "port", "2000")
	config.Store.Dir = configFile.MustValue("store", "dir")
	config.Store.PidFile = configFile.MustValue("store", "pidfile", "/tmp/kmdb-"+config.Listen.Port+".pid")
	config.Type.Primary = configFile.MustBool("type", "primary", true)
	if !config.Type.Primary {
		config.Type.SlaveOf = configFile.MustValue("type", "slaveOf")
	}
}
