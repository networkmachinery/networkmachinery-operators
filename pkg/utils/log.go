package utils

import (
	"os"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// LogErrAndExit logs the given error with msg and keysAndValues and calls `os.Exit(1)`.
func LogErrAndExit(err error, msg string, keysAndValues ...interface{}) {
	log.Log.Error(err, msg, keysAndValues...)
	os.Exit(1)
}

