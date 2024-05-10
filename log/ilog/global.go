package ilog

import (
	"log"
	"os"
	"sync"
)

//nolint:gochecknoglobals
var (
	_globalLogger   Logger = NewBuilder(DebugLevel, os.Stdout).Build() //nolint:revive
	_globalLoggerMu sync.RWMutex
)

func L() Logger { //nolint:ireturn
	_globalLoggerMu.RLock()
	l := _globalLogger
	_globalLoggerMu.RUnlock()
	return l
}

func Global() Logger { //nolint:ireturn
	return L()
}

func SetGlobal(logger Logger) (rollback func()) {
	_globalLoggerMu.Lock()
	backup := _globalLogger
	_globalLogger = logger
	_globalLoggerMu.Unlock()
	return func() {
		SetGlobal(backup)
	}
}

//nolint:gochecknoglobals
var stdLoggerMu sync.Mutex

func SetStdLogger(l Logger) (rollback func()) {
	stdLoggerMu.Lock()
	defer stdLoggerMu.Unlock()

	backupFlags := log.Flags()
	backupPrefix := log.Prefix()
	backupWriter := log.Writer()

	log.SetFlags(0)
	log.SetPrefix("")
	const skipForStgLogger = 2
	log.SetOutput(l.Copy().AddCallerSkip(skipForStgLogger))

	return func() {
		log.SetFlags(backupFlags)
		log.SetPrefix(backupPrefix)
		log.SetOutput(backupWriter)
	}
}
