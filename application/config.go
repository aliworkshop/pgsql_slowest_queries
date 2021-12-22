package application

import (
	"github.com/spf13/viper"
	"os"
)

type config struct {
	Sql     SqlConfig
	General struct {
		Address string
		Port    string
	}
}

func InitConfig(paths ...string) *viper.Viper {
	var path string
	if paths != nil && len(path) > 0 {
		path = paths[0]
	}
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
		if path == "" {
			path = "./config.yaml"
		}
	}
	v := viper.New()
	f, err := os.Open(path)
	if err != nil {
		panic("cannot read config: " + err.Error())
	}
	v.SetConfigType("yaml")
	err = v.ReadConfig(f)
	if err != nil {
		panic("cannot read config: " + err.Error())
	}
	return v
}
