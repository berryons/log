package log

import (
	"github.com/berryons/log/wrapper"
	"log"
)

var gLogger Logger
var gIsLogging bool

func init() {
	Enable()
	gLogger = wrapper.New("DEBUG", "", "", "", 3, true)
}

func SetLogger(logLevel, logPath, logFileName, fileLogLevel string, callerSkip int, isJson bool) {
	gLogger = wrapper.New(logLevel, logPath, logFileName, fileLogLevel, callerSkip, isJson)
}

func Enable() {
	gIsLogging = true
}

func Disable() {
	gIsLogging = false
}

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)

	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)

	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Print(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Print(v...)
		} else {
			log.Print(v...)
		}
	}
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Printf(format string, v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Printf(format, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Println(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Println(v...)
		} else {
			log.Println(v...)
		}
	}
}

// Fatal is equivalent to [Print] followed by a call to [os.Exit](1).
func Fatal(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Fatal(v...)
		} else {
			log.Fatal(v...)
		}
	}
}

// Fatalf is equivalent to [Printf] followed by a call to [os.Exit](1).
func Fatalf(format string, v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Fatalf(format, v...)
		} else {
			log.Fatalf(format, v...)
		}
	}
}

// Fatalln is equivalent to [Println] followed by a call to [os.Exit](1).
func Fatalln(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Fatalln(v...)
		} else {
			log.Fatalln(v...)
		}
	}
}

// Panic is equivalent to [Print] followed by a call to panic().
func Panic(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Panic(v...)
		} else {
			log.Panic(v...)
		}
	}
}

// Panicf is equivalent to [Printf] followed by a call to panic().
func Panicf(format string, v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Panicf(format, v...)
		} else {
			log.Println(v...)
		}
	}
}

// Panicln is equivalent to [Println] followed by a call to panic().
func Panicln(v ...any) {
	if gIsLogging {
		if gLogger != nil {
			gLogger.Panicln(v...)
		} else {
			log.Panicln(v...)
		}
	}
}
