package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	SupabaseURL        string
	SupabaseJWTSecret  string
	OpenAIAPIKey       string
	MarketDataAPIKey   string
	AllowedOrigin      string
}

func Load() Config {
	// .env is optional in production where real env vars are injected instead.
	_ = godotenv.Load()

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		SupabaseURL:       getEnv("SUPABASE_URL", ""),
		SupabaseJWTSecret: getEnv("SUPABASE_JWT_SECRET", ""),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		MarketDataAPIKey:  getEnv("MARKET_DATA_API_KEY", ""),
		AllowedOrigin:     getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
