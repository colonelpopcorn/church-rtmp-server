package main

import (
	"os"
	"runtime"
)

const SUCCESS_KEY = "success"
const RESPONSE_MESSAGE_KEY = "responseMessage"
const SQLITE_DATABASE = "sqlite-database.db"

func GetStreamingServerPath() string {
	osName := runtime.GOOS

	switch osName {
	case `windows`:
		return `C:\ProgramData\StreamingServer`
	case `darwin`:
		return `/Library/Applications/StreamingServer`
	default:
		optionalPath, ok := os.LookupEnv("STREAMING_SERVER_PATH")
		if ok {
			return optionalPath
		}
		return `/etc/streaming-server`
	}
}
