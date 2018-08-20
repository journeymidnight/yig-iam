package helper

import (
	"github.com/BurntSushi/toml"
	"os"
	"math/rand"
	"fmt"
	"reflect"
	"flag"
)

func Ternary(IF bool, THEN interface{}, ELSE interface{}) interface{} {
	if IF {
		return THEN
	} else {
		return ELSE
	}
}

// Static alphaNumeric table used for generating unique request ids
var alphaNumericTable = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

var NumericTable = []byte("0123456789")

func GenerateRandomId() []byte {
	alpha := make([]byte, 16, 16)
	for i := 0; i < 16; i++ {
		n := rand.Intn(len(alphaNumericTable))
		alpha[i] = alphaNumericTable[n]
	}
	return alpha
}

func GenerateRandomIdByLength(length int) []byte {
	alpha := make([]byte, length, length)
	for i := 0; i < length; i++ {
		n := rand.Intn(len(alphaNumericTable))
		alpha[i] = alphaNumericTable[n]
	}
	return alpha
}

func GenerateRandomNumberId() []byte {
	alpha := make([]byte, 16, 16)
	for i := 0; i < 16; i++ {
		n := rand.Intn(len(NumericTable))
		alpha[i] = NumericTable[n]
	}
	return alpha
}

func GenerateKey() ([]byte, []byte){
	accessKey := GenerateRandomIdByLength(20)
	accessSecret := GenerateRandomIdByLength(40)
	return accessKey, accessSecret
}

func GenerateUserId() string {
	return "u-" + string(GenerateRandomId())
}

func GenerateProjectId() string {
	return "p-" + string(GenerateRandomId())
}

var Config *Configuration

const CONFIGPATH = "/etc/yig-iam/conf.toml"

type Configuration struct {
	ManageKey                  string
	ManageSecret               string
	AccessLog                  string
	LogPath                    string
	LogLevel                   string
	PidFile                    string
	BindPort                   int
	RbacDataSource             string
	UserDataSource             string
	TokenExpire                int  //second
}

func SetupConfig() {
	Config, _ = LoadConfig()
	fmt.Println(Config)
	Logger = GetLog()
}

func DefaultConfiguration() *Configuration {
	cfg := &Configuration{
		ManageKey: "key",
		ManageSecret: "secret",
		AccessLog:  "/var/log/nier/access.log",
		LogPath:  "/var/log/nier/nier.log",
		LogLevel: "info",
		RbacDataSource:   "root:@tcp(127.0.0.1:3306)/",
		UserDataSource:   "root:@tcp(127.0.0.1:3306)/",
		TokenExpire:  60*60*10,
	}
	return cfg
}

func LoadConfig() (*Configuration, error) {
	rtConfig := DefaultConfiguration()
	if _, err := os.Stat(CONFIGPATH); err != nil {
		fmt.Fprintln(os.Stderr,"config file does exsit,skipped config file")
	} else {
		_, err = toml.DecodeFile("/etc/yig-iam/conf.toml", &rtConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr,"failed to decode config file,skipped config file", err)
		}
	}
	mergeConfig(rtConfig, configFromFlag())
	return rtConfig, nil
}

func configFromFlag() *Configuration {
	cfg := &Configuration{}
	flag.StringVar(&cfg.AccessLog, "accesslog", "", "path for access file")
	flag.StringVar(&cfg.LogPath, "logpath", "", "path for the log file")
	flag.StringVar(&cfg.LogLevel, "loglevel", "", "using standard go library")
	flag.StringVar(&cfg.RbacDataSource, "rbacdatasource", "", "using standard mysql datasource")
	flag.StringVar(&cfg.UserDataSource, "userdatasource", "", "using standard mysql datasource")
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

