package helper

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	alpha := make([]byte, 16, 16)
	for i := 0; i < 16; i++ {
		n := rand.Intn(len(alphaNumericTable))
		alpha[i] = alphaNumericTable[n]
	}
	return alpha
}

func GenerateRandomIdByLength(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	alpha := make([]byte, length, length)
	for i := 0; i < length; i++ {
		n := rand.Intn(len(alphaNumericTable))
		alpha[i] = alphaNumericTable[n]
	}
	return alpha
}

func GenerateRandomNumberId() []byte {
	rand.Seed(time.Now().UnixNano())
	alpha := make([]byte, 16, 16)
	for i := 0; i < 16; i++ {
		n := rand.Intn(len(NumericTable))
		alpha[i] = NumericTable[n]
	}
	return alpha
}

type Config struct {
	ManageKey                      string
	ManageSecret                   string
	LogPath                        string
	PanicLogPath                   string
	PidFile                        string
	BindPort                       int
	DatabaseConnectionString       string
	CasbinDbString                 string
	DebugMode                      bool
	LogLevel                       int //1-20
	TokenExpire                    int //second
	AccessKey, SecretKey, S3Domain string
}

type config struct {
	ManageKey                      string
	ManageSecret                   string
	LogPath                        string
	PanicLogPath                   string
	PidFile                        string
	BindPort                       int
	DatabaseConnectionString       string
	CasbinDbString                 string
	DebugMode                      bool
	LogLevel                       int //1-20
	TokenExpire                    int //second
	AccessKey, SecretKey, S3Domain string
}

var CONFIG Config

func SetupConfig() {
	f, err := os.Open("/etc/yig/iam.json")
	if err != nil {
		panic("Cannot open iam.json")
	}
	defer f.Close()

	var c config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		panic("Failed to parse yig.json: " + err.Error())
	}

	// setup CONFIG with defaults
	CONFIG.ManageKey = c.ManageKey
	CONFIG.ManageSecret = c.ManageSecret
	CONFIG.LogPath = c.LogPath
	CONFIG.PanicLogPath = c.PanicLogPath
	CONFIG.PidFile = c.PidFile
	CONFIG.BindPort = c.BindPort
	CONFIG.DatabaseConnectionString = c.DatabaseConnectionString
	CONFIG.CasbinDbString = c.CasbinDbString
	CONFIG.DebugMode = c.DebugMode
	CONFIG.LogLevel = Ternary(c.LogLevel == 0, 5, c.LogLevel).(int)
	CONFIG.TokenExpire = Ternary(c.TokenExpire == 0, 28800, c.TokenExpire).(int)
	CONFIG.AccessKey = c.AccessKey
	CONFIG.SecretKey = c.SecretKey
	CONFIG.S3Domain = c.S3Domain
}
