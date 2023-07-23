package log

import (
	"fmt"
	"io"
	"log"
	stdLog "log"
	"os"
)

var (
	infoLogger  *stdLog.Logger
	errLogger   *stdLog.Logger
	panicLogger *stdLog.Logger
)

func Init() {
	setLogger(os.Stderr, log.Ldate|log.Ltime|log.Lshortfile)
}

func setLogger(out io.Writer, flag int) {
	infoLogger = stdLog.New(out, "[INFO]\t", flag)
	errLogger = stdLog.New(out, "[ERROR]\t", flag)
	panicLogger = stdLog.New(out, "[PANIC]\t", flag)
}

func Info(v ...any) {
	infoLogger.Print(v...)
}

func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Error(v ...any) {
	errLogger.Print(v...)
}

func Errorf(format string, v ...any) {
	errLogger.Printf(format, v...)
}

func Panic(v ...any) {
	panicLogger.Print(v...)
}

func Panicf(format string, v ...any) {
	panicLogger.Printf(format, v...)
	panic(fmt.Sprintf(format, v...))
}
