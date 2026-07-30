package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"finedu-backend/internal/config"
	"finedu-backend/internal/db"
	"finedu-backend/internal/middleware"
	"finedu-backend/internal/routes"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	jwks, err := middleware.NewSupabaseKeyfunc(ctx, cfg.SupabaseURL)
	if err != nil {
		log.Fatalf("failed to load Supabase JWKS: %v", err)
	}

	r := gin.Default()
	// No reverse proxy in front of this server today, so trust none: without
	// this, Gin's default trusts every proxy and honors a client-supplied
	// X-Forwarded-For header, letting anyone spoof the IP RateLimit keys on.
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to configure trusted proxies: %v", err)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.Register(r, pool, jwks, cfg)

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
