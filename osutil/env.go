package osutil

import (
	"os"
	"strconv"
)

func GetEnvStr(key string, defaultValue ...string) string {
	envStr := os.Getenv(key)

	if envStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return ""
		}
	}

	return envStr
}

func GetEnvStrReq(key string) string {
	envStr := GetEnvStr(key)

	if envStr == "" {
		log.Fatalf("Missing required environment variable '%s'", key)
	}

	return envStr
}

func GetEnvInt(key string, defaultValue ...int) int {
	envStr := os.Getenv(key)

	envInt, err := strconv.Atoi(envStr)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return 0
		}
	}

	return envInt
}

func GetEnvIntReq(key string) int {
	envStr := os.Getenv(key)

	envInt, err := strconv.Atoi(envStr)
	if err != nil {
		log.Fatalf("Missing required environment variable '%s'", key)
	}

	return envInt
}
