package helper

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

var Config *Configuration

const CONFIGPATH = "/etc/nier/conf.toml"

type Configuration struct {
	Accesslog string
	Logpath  string
	Loglevel string
	RbacDataSource   string
	UserDataSource   string
	Prometheus string
	Catkeeper string
	Serveblk string
	TokenExpire  int
}

func SetupConfig() {
	Config, _ = LoadConfig()
	fmt.Println(Config)
	Logger = GetLog()
}

func DefaultConfiguration() *Configuration {
	cfg := &Configuration{
		Accesslog:  "/var/log/nier/access.log",
		Logpath:  "/var/log/nier/nier.log",
		Loglevel: "info",
		RbacDataSource:   "root:@tcp(10.72.84.145:4000)/",
		UserDataSource:   "root:@tcp(10.72.84.145:4000)/iam",
		Prometheus: "127.0.0.1:9090",
		Catkeeper: "127.0.0.1:9091",
		Serveblk: "127.0.0.1:18080",
		TokenExpire:  60*60*10,
	}
	return cfg
}

func LoadConfig() (*Configuration, error) {
	rtConfig := DefaultConfiguration()
	if _, err := os.Stat(CONFIGPATH); err != nil {
		fmt.Fprintln(os.Stderr,"config file does exsit,skipped config file")
	} else {
		_, err = toml.DecodeFile("/etc/nier/conf.toml", &rtConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr,"failed to decode config file,skipped config file", err)
		}
	}
	mergeConfig(rtConfig, configFromFlag())
	return rtConfig, nil
}

func configFromFlag() *Configuration {
	cfg := &Configuration{}
	flag.StringVar(&cfg.Accesslog, "accesslog", "", "path for access file")
	flag.StringVar(&cfg.Logpath, "logpath", "", "path for the log file")
	flag.StringVar(&cfg.Loglevel, "loglevel", "", "using standard go library")
	flag.StringVar(&cfg.RbacDataSource, "rbacdatasource", "", "using standard mysql datasource")
	flag.StringVar(&cfg.UserDataSource, "userdatasource", "", "using standard mysql datasource")
	flag.StringVar(&cfg.Prometheus, "prometheus", "", "prometheus endpoint ip:port")
	flag.StringVar(&cfg.Catkeeper, "catkeeper", "", "catkeeper endpoint ip:port")
	flag.StringVar(&cfg.Serveblk, "serveblk", "", "serveblk endpoint ip:port")
	flag.Parse()
	return cfg
}

func mergeConfig(defaultcfg, filecfg interface{}) {
	v1 := reflect.ValueOf(filecfg).Elem()
	v := reflect.ValueOf(defaultcfg).Elem()
	mergeValue(v, v1)
}

func mergeValue(v, v1 reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			if v.Field(i).CanSet() && !v1.Field(i).IsNil() {
				mergeValue(v.Field(i).Elem(), v1.Field(i).Elem())
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		case reflect.Bool:
			if v.Field(i).CanSet() {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		case reflect.Int:
			if v.Field(i).CanSet() && v1.Field(i).Int() != 0 {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		default:
			if v.Field(i).CanSet() && v1.Field(i).Len() != 0 {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		}
	}
}
