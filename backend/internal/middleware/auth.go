package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextUserIDKey is the gin context key holding the authenticated user's Supabase ID (the "sub" claim).
const ContextUserIDKey = "userID"

// NewSupabaseKeyfunc fetches Supabase's public JWT signing keys (JWKS) and
// keeps them refreshed in the background, so tokens can be verified without
// the backend ever holding a shared secret. This also means Supabase can
// rotate its signing key (old + new keys coexist during rollover) without
// any backend redeploy.
func NewSupabaseKeyfunc(ctx context.Context, supabaseURL string) (keyfunc.Keyfunc, error) {
	jwksURL := strings.TrimRight(supabaseURL, "/") + "/auth/v1/.well-known/jwks.json"
	return keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
}

// RequireAuth validates the Supabase-issued JWT on the Authorization header
// against Supabase's published JWKS and rejects the request if it's missing,
// malformed, expired, or signed with an untrusted key/algorithm.
func RequireAuth(k keyfunc.Keyfunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be a bearer token"})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, k.Keyfunc, jwt.WithValidMethods([]string{"ES256", "RS256"}))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		sub, err := claims.GetSubject()
		if err != nil || sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing subject claim"})
			return
		}

		c.Set(ContextUserIDKey, sub)
		c.Next()
	}
}
