package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type Log struct {
	// dont forget
	Logging        *zap.Logger
	Sugar          *zap.SugaredLogger
	CurrentLogTime string
	Level          string
}

var LOG *Log

func InitDefaultLogger(level string, serviceName string) {
	LOG = &Log{}
	LOG = loggerInit("", level, serviceName)
	go loggerConfigRotate()
	DefaultConfigManager.HotUpdateTarget = append(DefaultConfigManager.HotUpdateTarget,LOG)
}

// InitNewLogger 传入保存日志目录的绝对路径
func InitNewLogger(logPath string, level string, serviceName string) {
	LOG = &Log{}
	LOG = loggerInit(logPath, level, serviceName)
	go loggerConfigRotate()
	DefaultConfigManager.HotUpdateTarget = append(DefaultConfigManager.HotUpdateTarget,LOG)
}

func loggerInit(logPath string, logLevel string, serviceName string) *Log {
	config := getZapConfig(logPath, logLevel, serviceName)
	initLogger, err := config.Build()
	if err != nil {
		panic("log init fail")
	}
	//initLogger.Info("create succuess.")
	initSugar := initLogger.Sugar()
	//initSugar.Info("create sugar success.")
	logger := &Log{}
	logger.Logging = initLogger
	logger.Sugar = initSugar
	logger.CurrentLogTime = time.Now().Format("2006-01-02")
	return logger
}

// loggerConfigRotate rebuild zap.logger everyday
func loggerConfigRotate() {
	duration := time.Duration(time.Second * 60)
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			if time.Now().Format("2006-01-02") != LOG.CurrentLogTime {
				newLogger := loggerInit("", LOG.Level, "rpc-test")
				// ?
				LOG = newLogger
			}
			ticker = time.NewTicker(duration)
		}
	}
}

func initLoggingFile(path string) {
	index := strings.LastIndex(path, "/")
	pathDir := path[:index]
	err := createFile(pathDir)
	if err != nil {
		panic("create log file dir error")
	}
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(path)
		defer func() {
			_ = file.Close()
		}()
		if err != nil {
			panic("create logFile fail")
		}
	}
}

//调用os.MkdirAll递归创建文件夹
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func (l *Log) Sync() {
	_ = l.Logging.Sync()
}

func (l *Log) Debugf(template string, args ...interface{}) {
	l.Sugar.Debugf(template, args ...)
}

func (l *Log) Infof(template string, args ...interface{}) {
	l.Sugar.Infof(template, args ...)
}

func (l *Log) Errorf(template string, args ...interface{}) {
	l.Sugar.Errorf(template, args ...)
}

func (l *Log) Warnf(template string, args ...interface{}) {
	l.Sugar.Warnf(template, args ...)
}

func (l *Log) Fatalf(template string, args ...interface{}) {
	l.Sugar.Fatalf(template, args ...)
}

func (l *Log) Debug(log string) {
	l.Sugar.Debug(log)
}

func (l *Log) Info(log string) {
	l.Sugar.Info(log)
}
func (l *Log) Warn(log string) {
	l.Sugar.Warn(log)
}
func (l *Log) Error(log string) {
	l.Sugar.Error(log)
}

func (l *Log) Fatal(log string) {
	l.Sugar.Fatal(log)
}

func getZapConfig(pathInfo string, level string, serviceName string) zap.Config {
	var logPath string
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	if pathInfo != "" {
		logPath = pathInfo
	} else {
		dirPath, err := os.Getwd()
		if err != nil {
			panic("Get project dir error.")
		}
		logPath = dirPath + "/logs"
		// logPath = "./logs/" + time.Now().Format("2006-01-02") + "_" + serviceName + ".log"
	}
	logPath += "/" + time.Now().Format("2006-01-02") + "_" + serviceName + ".log"
	initLoggingFile(logPath)

	var Level zap.AtomicLevel
	switch level {
	case "Debug":
		Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "Info":
		Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "Error":
		Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "Warn":
		Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "Fatal":
		Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	}

	defaultConfig := zap.Config{
		Level:            Level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		InitialFields:    map[string]interface{}{"serviceName": serviceName},
		OutputPaths:      []string{"stdout", logPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	return defaultConfig
}

func (l *Log)UpdateConfig()  {
	// update config like log level,path
	LOG = loggerInit("", DefaultToml.Log.Level, DefaultToml.Server.ServiceName)

}
