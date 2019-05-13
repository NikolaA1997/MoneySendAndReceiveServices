package util

import (
	"log"
	"runtime"
	"strings"
)

func ErrorLog(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func MessageFormat(message string) []string{
	parts := strings.Split(strings.Trim(message, "{}")," ")
	return parts
}
