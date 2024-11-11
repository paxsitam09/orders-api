package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	ServerPort    uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress:  "redis-10670.c340.ap-northeast-2-1.ec2.redns.redis-cloud.com:10670",
		RedisPassword: "ocknMiTUWBO52WZRaNBfm8lA2ByPrn5s",
		RedisDB:       0,
		ServerPort:    3000,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}
