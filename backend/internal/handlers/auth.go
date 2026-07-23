package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/config"
)

type signUpRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
	FullName     string `json:"full_name" binding:"required"`
	AgreeToTerms bool   `json:"agree_to_terms"`
}

// SignUp validates the request, then creates the user by calling Supabase's
// own public signup endpoint server-side. This mirrors exactly what the
// Supabase JS client does from the browser (same confirmation-email
// behavior), but means our validation rules can't be skipped by anyone
// going through our app.
// Owner: Auth
func SignUp(cfg config.Config) gin.HandlerFunc {
	client := &http.Client{}

	return func(c *gin.Context) {
		var req signUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if issues := passwordIssues(req.Password); len(issues) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must have: " + strings.Join(issues, ", ")})
			return
		}

		if !isValidFullName(req.FullName) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "full name must be 2-100 characters and contain only letters, spaces, hyphens, apostrophes, or periods"})
			return
		}

		if !req.AgreeToTerms {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you must agree to the terms of service and privacy policy"})
			return
		}

		if err := callSupabaseSignUp(c.Request.Context(), client, cfg, req); err != nil {
			status, message := normalizeSignUpError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "check your email to confirm your account"})
	}
}

func passwordIssues(password string) []string {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	var issues []string
	if len(password) < 8 {
		issues = append(issues, "at least 8 characters")
	}
	if !hasLower {
		issues = append(issues, "a lowercase letter")
	}
	if !hasUpper {
		issues = append(issues, "an uppercase letter")
	}
	if !hasDigit {
		issues = append(issues, "a number")
	}
	if !hasSpecial {
		issues = append(issues, "a special character")
	}
	return issues
}

func isValidFullName(name string) bool {
	trimmed := strings.TrimSpace(name)
	// Rune count, not len() (byte count), so multi-byte characters (accented
	// letters, Cyrillic, CJK, etc.) are bounded the same way as the
	// frontend's JS .length check.
	runeCount := utf8.RuneCountInString(trimmed)
	if runeCount < 2 || runeCount > 100 {
		return false
	}
	for _, r := range trimmed {
		if !unicode.IsLetter(r) && r != ' ' && r != '\'' && r != '-' && r != '.' {
			return false
		}
	}
	return true
}

type supabaseSignUpError struct {
	statusCode int
	errorCode  string
	message    string
}

func (e *supabaseSignUpError) Error() string {
	return e.message
}

func callSupabaseSignUp(ctx context.Context, client *http.Client, cfg config.Config, req signUpRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	payload, err := json.Marshal(gin.H{
		"email":    req.Email,
		"password": req.Password,
		"data": gin.H{
			"full_name":      req.FullName,
			"terms_accepted": req.AgreeToTerms,
		},
	})
	if err != nil {
		return &supabaseSignUpError{statusCode: http.StatusInternalServerError, message: err.Error()}
	}

	signupURL := strings.TrimRight(cfg.SupabaseURL, "/") + "/auth/v1/signup"
	if cfg.AllowedOrigin != "" {
		signupURL += "?redirect_to=" + url.QueryEscape(cfg.AllowedOrigin)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, signupURL, bytes.NewReader(payload))
	if err != nil {
		return &supabaseSignUpError{statusCode: http.StatusInternalServerError, message: err.Error()}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", cfg.SupabaseAnonKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return &supabaseSignUpError{statusCode: http.StatusBadGateway, message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	var body struct {
		ErrorCode        string `json:"error_code"`
		Msg              string `json:"msg"`
		ErrorDescription string `json:"error_description"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	message := body.Msg
	if message == "" {
		message = body.ErrorDescription
	}

	return &supabaseSignUpError{statusCode: resp.StatusCode, errorCode: body.ErrorCode, message: message}
}

func normalizeSignUpError(err error) (int, string) {
	supaErr, ok := err.(*supabaseSignUpError)
	if !ok {
		log.Printf("signup: unexpected error: %v", err)
		return http.StatusInternalServerError, "unable to create account, please try again"
	}

	log.Printf("signup: supabase rejected signup (status=%d code=%s): %s", supaErr.statusCode, supaErr.errorCode, supaErr.message)

	lowerMsg := strings.ToLower(supaErr.message)
	switch {
	case supaErr.errorCode == "user_already_exists" || strings.Contains(lowerMsg, "already registered"):
		return http.StatusConflict, "an account with this email already exists"
	case supaErr.errorCode == "weak_password" || strings.Contains(lowerMsg, "password"):
		return http.StatusBadRequest, "password does not meet requirements"
	case supaErr.statusCode >= http.StatusInternalServerError || supaErr.statusCode == http.StatusBadGateway:
		return http.StatusBadGateway, "unable to reach the authentication service, please try again"
	default:
		return http.StatusBadRequest, "unable to create account, please try again"
	}
}
