package wrapper

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func New(logLevel, logPath, logFileName, fileLogLevel string, callerSkip int, isJson bool) *Wrapper {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	realLogPath := path.Join(logPath, logFileName)
	var cores []zapcore.Core

	// Create Console Logger
	newCore, err := NewConsoleLog(convertZapLogLevelFromString(logLevel), isJson).NewCore()
	if err != nil {
		log.Fatal(err)
	}
	cores = append(cores, newCore)

	// Create File Logger
	if len(logPath) != 0 {
		fileLogInfo := NewFileLog(logPath, logFileName, convertZapLogLevelFromString(fileLogLevel), isJson)
		newCore, err := fileLogInfo.NewCore()
		if err != nil {
			log.Fatal(err)
		}
		realLogPath = fileLogInfo.LogPath
		cores = append(cores, newCore)
	}

	// Combine logger cores
	core := zapcore.NewTee(cores...)

	// Create the logger with additional context information (caller, stack trace)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip))

	return &Wrapper{Logger: logger, LogPath: realLogPath}
}

type Wrapper struct {
	Logger  *zap.Logger
	LogPath string
}

func (iSelf Wrapper) GetLogger() interface{} {
	return iSelf.Logger
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func (iSelf Wrapper) Print(v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
		iSelf.Logger.Info(checkedStr)
		iSelf.sync()
	}
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func (iSelf Wrapper) Printf(format string, v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf(format, v...))
		iSelf.Logger.Info(checkedStr)
		iSelf.sync()
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func (iSelf Wrapper) Println(v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
		iSelf.Logger.Info(checkedStr)
		iSelf.sync()
	}
}

// Fatal is equivalent to [Print] followed by a call to [os.Exit](1).
func (iSelf Wrapper) Fatal(v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
		iSelf.Logger.Fatal(checkedStr)
		iSelf.sync()
	}
	os.Exit(1)
}

// Fatalf is equivalent to [Printf] followed by a call to [os.Exit](1).
func (iSelf Wrapper) Fatalf(format string, v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf(format, v...))
		iSelf.Logger.Fatal(checkedStr)
		iSelf.sync()
	}
	os.Exit(1)
}

// Fatalln is equivalent to [Println] followed by a call to [os.Exit](1).
func (iSelf Wrapper) Fatalln(v ...any) {
	if iSelf.Logger != nil {
		checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
		iSelf.Logger.Fatal(checkedStr)
		iSelf.sync()
	}
	os.Exit(1)
}

// Panic is equivalent to [Print] followed by a call to panic().
func (iSelf Wrapper) Panic(v ...any) {
	checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
	if iSelf.Logger == nil {
		iSelf.Logger.Panic(checkedStr)
		iSelf.sync()
	}
	panic(checkedStr)
}

// Panicf is equivalent to [Printf] followed by a call to panic().
func (iSelf Wrapper) Panicf(format string, v ...any) {
	checkedStr := checkNewLine(fmt.Sprintf(format, v...))
	if iSelf.Logger == nil {
		iSelf.Logger.Panic(checkedStr)
		iSelf.sync()
	}
	panic(checkedStr)
}

// Panicln is equivalent to [Println] followed by a call to panic().
func (iSelf Wrapper) Panicln(v ...any) {
	checkedStr := checkNewLine(fmt.Sprintf("%s", v...))
	if iSelf.Logger == nil {
		iSelf.Logger.Panic(checkedStr)
		iSelf.sync()
	}
	panic(checkedStr)
}

func (iSelf Wrapper) sync() {
	err := iSelf.Logger.Sync()
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] Logging error: %s", err)
		if iSelf.Logger != nil {
			iSelf.Logger.Error(errMsg)
		} else {
			log.Println(errMsg)
		}
	}
}

func checkNewLine(str string) string {
	checked := str
	if index := strings.LastIndex(checked, "\n"); index != -1 {
		checked = checked[:index]
	}
	return checked
}

func convertZapLogLevelFromString(logLevelStr string) zapcore.Level {
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "DPANIC":
		return zapcore.DPanicLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InvalidLevel
	}
}
