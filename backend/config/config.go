package config

import "os"

var (
	DBUrl     = os.Getenv("DATABASE_URL")
	RedisUrl  = os.Getenv("REDIS_URL")
	JwtSecret = os.Getenv("JWT_SECRET")
)
