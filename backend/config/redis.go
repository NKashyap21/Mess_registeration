package config

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadRedisConfig() *RedisConfig {
	dbStr := getEnv("REDIS_DB", "0")
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Printf("Invalid REDIS_DB value, using default: %v", err)
		db = 0
	}

	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

func ConnectRedis() {
	config := LoadRedisConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connected successfully!")
}

func GetRedisClient() *redis.Client {
	return RedisClient
}
