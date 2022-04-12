package helper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Go starts go-routine in a safe way.
// see protect function in https://golang.org/ref/spec#Handling_panics
func Go(name string, g func()) {
	go func() {
		defer RecoverAsLogError(name)
		g()
	}()
}

// ResetEnv sets the environment variables to the list of variables pulled from `os.Environ()`
func ResetEnv(environ []string) (err error) {
	os.Clearenv()
	for _, env := range environ {
		keyVal := strings.SplitN(env, "=", 2)
		err = os.Setenv(keyVal[0], keyVal[1])
	}

	return err
}

func zlogger() *zap.Logger {
	res, err2 := zap.NewProduction()
	ErrorHandler(err2, "zap", "logger initialization")
	return res
}

// PanicOnError panics if error is not nil
// It's recommended to supplement additional context information using
// `errors.Wrapf(err, "info %s", data)`
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// ErrorHandler is a universal handler that panics unless err is nil
func ErrorHandler(err error, namespace, msg string) {
	if err != nil {
		panic(errors.Wrapf(err, "ERROR in '%s': %s\n", namespace, msg))
	}
}

// ErrorHandlerf is a universal handler that panics unless err is nil
func ErrorHandlerf(err error, namespace, format string, args ...interface{}) {
	if err != nil {
		ErrorHandler(err, namespace, fmt.Sprintf(format, args...))
	}
}

// RecoverAsFalse recovers from panic and returns false
func RecoverAsFalse(name string, res *bool) {
	if err2 := recover(); err2 != nil {
		switch err2 := err2.(type) {
		case error:
			zlogger().Error(name, zap.Error(err2))
		default:
		}
		log.Printf("IGNORING ERROR in %s (returning false): %+v", name, err2)
		*res = false
	}
}

// RecoverAsLogError is a universal error recovery that logs error
func RecoverAsLogError(label string) {
	if err2 := recover(); err2 != nil {
		logError(label, err2)
	}
}

func logError(label string, err2 interface{}) {
	switch err2 := err2.(type) {
	case error:
		zlogger().Error(label, zap.Error(err2))
	default:
	}
	log.Printf("IGNORING ERROR in %s: %+v\n", label, err2)
}

// RecoverAsLogErrorf is a universal error recovery that logs error
func RecoverAsLogErrorf(format string, args ...interface{}) {
	// RecoverAsLogError(label) - this doesn't work due to limitation in `recover()`
	// in particular, it works only in the first level of nesting.
	// See https://golang.org/ref/spec#Handling_panics and
	// https://stackoverflow.com/questions/49344478
	if err2 := recover(); err2 != nil {
		label := fmt.Sprintf(format, args...)
		logError(label, err2)
	}
}

// RecoverToErrorVar recovers and places the recovered error into the given variable
func RecoverToErrorVar(name string, err *error) {
	err2 := recover()
	if err2 != nil {
		log.Printf("RecoverToErrorVar2 (%s) (err=%+v), (err2: %+v\n", name, *err, err2)
		switch err2 := err2.(type) {
		case error:
			err4 := errors.Wrapf(err2, "%s: Recover from panic", name)
			*err = err4
		case string:
			err4 := errors.New(name + ": Recover from string-panic: " + err2)
			*err = err4
		default:
			err4 := errors.New(fmt.Sprintf("%s: Recover from unknown-panic: %+v", name, err2))
			*err = err4
		}
	}
}

func WriteToCSVFile(filename string, header []string, data []map[string]interface{}) (int, error) {
	handle, err := os.OpenFile(filename, os.O_CREATE, os.ModeAppend)
	if err != nil {
		return 0, err
	}

	rows := 0
	for _, v := range data {

		lrows, err2 := WriteToCSV(handle, header, v)
		rows += lrows
		if err2 != nil {
			return rows, err2
		}
	}

	return rows, nil
}

func WriteToCSV(handle *os.File, header []string, data map[string]interface{}) (int, error) {
	rows := 0

	for _, v := range header {
		_, err := fmt.Fprintf(handle, "%s", v)
		if err != nil {
			return rows, err
		}
	}

	for _, v := range data {
		rows++

		_, err := fmt.Fprintf(handle, "%+v", v)
		if err != nil {
			return rows, err
		}
	}

	return rows, nil
}
