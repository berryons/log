package wrapper

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func NewConsoleLog(logLevel zapcore.Level, isJson bool) *ConsoleLog {
	consoleLogging := new(ConsoleLog)
	consoleLogging.ILog = interface{}(consoleLogging).(ILog)
	consoleLogging.LogLevel = logLevel
	consoleLogging.IsJson = isJson
	return consoleLogging
}

type ILog interface {
	NewCore() (zapcore.Core, error)
	GetEncoder() zapcore.Encoder
	GetWriter() zapcore.WriteSyncer
}

type ConsoleLog struct {
	ILog
	LogLevel zapcore.Level
	IsJson   bool
}

func (pSelf *ConsoleLog) NewCore() (zapcore.Core, error) {
	// Encoder
	encoder := pSelf.ILog.GetEncoder()

	// Writer
	writer := pSelf.ILog.GetWriter()

	// Create core
	return zapcore.NewCore(encoder, writer, pSelf.LogLevel), nil
}

func (pSelf *ConsoleLog) GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	if pSelf.IsJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (pSelf *ConsoleLog) GetWriter() zapcore.WriteSyncer {
	// Create syncer for console
	return zapcore.Lock(WrappedStdOutWriteSyncer{os.Stdout})
}

type WrappedStdOutWriteSyncer struct {
	// 윈도우에서는 괜찮으나, OSX, 리눅스에서 에러가 남.
	// OS X: inappropriate ioctl for device
	// Linux: sync /dev/stderr: invalid argument
	// os.Stdout handle 을 wrap 하여 Consol Sync 작업 시, 아무 일도 하지 않도록 함.
	stdOut *os.File
}

func (ws WrappedStdOutWriteSyncer) Write(p []byte) (n int, err error) {
	return ws.stdOut.Write(p)
}
func (ws WrappedStdOutWriteSyncer) Sync() error {
	return nil
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.234"
	wrappedContents := wrappingLogContents(t.Format(layout))
	enc.AppendString(wrappedContents)
}

func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[")
	zapcore.CapitalColorLevelEncoder(l, enc)
	enc.AppendString("]")
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	//wrappedContents := wrappingLogContents(caller.String()) // full
	wrappedContents := wrappingLogContents(caller.TrimmedPath()) // short

	enc.AppendString(wrappedContents)
}

func wrappingLogContents(contents string) string {
	return fmt.Sprintf("[%s]", contents)
}
