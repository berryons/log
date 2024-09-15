package wrapper

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewFileLog(logPath, logFileName string, logLevel zapcore.Level, isJson bool) *FileLog {
	fileLogging := new(FileLog)
	fileLogging.ILog = interface{}(fileLogging).(ILog)
	fileLogging.LogPath = logPath
	fileLogging.LogLevel = logLevel
	fileLogging.LogFileName = logFileName
	fileLogging.IsJson = isJson
	return fileLogging
}

type FileLog struct {
	ConsoleLog
	LogPath     string
	LogFileName string
}

// GetEncoder Override
func (pSelf *FileLog) GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	if pSelf.IsJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customFileLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetWriter Override
func (pSelf *FileLog) GetWriter() zapcore.WriteSyncer {
	pSelf.LogPath = checkLogPath(pSelf.LogPath, pSelf.LogFileName)

	// Set up lumberjack as a logger:
	fileWriter := &lumberjack.Logger{
		Filename:   pSelf.LogPath, // Or any other path
		MaxSize:    10,            // MB; after this size, a new log file is created
		MaxBackups: 100,           // Number of backups to keep
		MaxAge:     1,             // Days
		Compress:   true,          // Compress the backups using gzip
	}

	// Create syncer for writing to the file
	return zapcore.AddSync(fileWriter)
}

func customFileLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	wrappedContents := wrappingLogContents(l.CapitalString())
	enc.AppendString(wrappedContents)
}

func checkLogPath(logPath, logFileName string) string {
	logPathToSet := "./"

	if len(logPath) > 0 {
		logPathToSet = logPath
	} else {
		logRootPath, err := os.Getwd()
		if err == nil {
			logPathToSet = filepath.Dir(logRootPath)
		}
	}

	return filepath.Join(logPathToSet, checkLogFileName(logFileName))
}

func checkLogFileName(logFileName string) string {
	var response string

	if len(logFileName) == 0 {
		response = fmt.Sprintf("%s.log", time.Now().Format(time.DateTime))
	} else {
		response = fmt.Sprintf("%s_%s.log", time.Now().Format(time.DateTime), logFileName)
	}

	response = strings.Replace(response, " ", "_", -1)
	response = strings.Replace(response, "-", "", -1)
	return strings.Replace(response, ":", "", -1)
}
