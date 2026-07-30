package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	SupabaseURL        string
	SupabaseAnonKey    string
	GeminiAPIKey       string
	AlphaVantageAPIKey string
	FinnhubAPIKey      string
	AllowedOrigin      string
}

func Load() Config {
	// .env is optional in production where real env vars are injected instead.
	_ = godotenv.Load()

	return Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		GeminiAPIKey:       getEnv("GEMINI_API_KEY", ""),
		AlphaVantageAPIKey: getEnv("ALPHA_VANTAGE_API_KEY", ""),
		FinnhubAPIKey:      getEnv("FINNHUB_API_KEY", ""),
		AllowedOrigin:      getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
