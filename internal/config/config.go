package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	S3          S3
}

type S3 struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

func Load() Config {
	return Config{
		Port:        getenv("PORT", "3000"),
		DatabaseURL: getenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5433/ecomm?sslmode=disable"),
		S3: S3{
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Region:    os.Getenv("S3_REGION"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			Bucket:    os.Getenv("S3_BUCKET"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
