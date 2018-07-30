package helper

import  (
	"os"
	"io"
	"path/filepath"
	"fmt"
	"errors"
	olog "github.com/sirupsen/logrus"
)
var Logger *olog.Logger

var level map[string]olog.Level = map[string]olog.Level{
	"info" : olog.InfoLevel,
	"warn" : olog.WarnLevel,
	"debug" : olog.DebugLevel,
	"error" : olog.ErrorLevel,
}

func OpenAccessLogFile () (*os.File, error) {
	if Config.Accesslog == "" {
		return nil, errors.New("No access log provided")
	}
	filePath := Config.Accesslog
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log path %s: %s", dir, err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %s", filePath, err)
	}

	return file, nil
}

func GetLog() *olog.Logger {
	var logdst io.Writer
	if Config.Logpath != "" {
		logdst, _ = os.OpenFile(Config.Logpath,os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	} else {
		logdst = os.Stdout
	}
	
	loglevel := Config.Loglevel
	Logger = olog.New()
	Logger.Out = logdst
	if _,exist := level[loglevel];exist {
		Logger.SetLevel(level[loglevel])
	}else{
		panic("wrong default log level")
	}
	return Logger
}
