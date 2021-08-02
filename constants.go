package main

import (
	"os"
	"runtime"
)

const SUCCESS_KEY = "success"
const RESPONSE_MESSAGE_KEY = "responseMessage"
const SQLITE_DATABASE = "sqlite-database.db"

func GetHomeFolder() string {
	osName := runtime.GOOS
	optionalPath, ok := os.LookupEnv("STREAMING_SERVER_PATH")

	if ok {
		return optionalPath
	}

	switch osName {
	case `windows`:
		return `C:\ProgramData\StreamingServer\`
	case `darwin`:
		return `/Library/Applications/StreamingServer/`
	default:
		return `/etc/streaming-server/`
	}
}
