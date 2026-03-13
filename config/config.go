package config

import (
	"os"
	"strconv"
)

type Config struct {
	ListenAddr   string
	TargetAddr   string
	SlowQueryMs  int
	N1WindowSecs int
	N1Threshold  int
}

func LoadConfig() Config {
	return Config{
		ListenAddr:   getEnv("SQLENS_LISTEN_ADDR", ":5433"),
		TargetAddr:   getEnv("SQLENS_TARGET_ADDR", "localhost:5432"),
		SlowQueryMs:  getEnvAsInt("SQLENS_SLOW_QUERY_MS", 100),
		N1WindowSecs: getEnvAsInt("SQLENS_N1_WINDOW_SECS", 10),
		N1Threshold:  getEnvAsInt("SQLENS_N1_THRESHOLD", 5),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
