package config

import "os"

type Config struct {
	Addr      string
	DBDSN     string
	UploadDir string
}

func Load() Config {
	return Config{
		Addr:      getEnv("APP_ADDR", ":8080"),
		DBDSN:     getEnv("DB_DSN", "ctf:ctf@tcp(127.0.0.1:3306)/campus_ctf?charset=utf8mb4&parseTime=True&loc=Local"),
		UploadDir: getEnv("APP_UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

