package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/musiermoore/ticketing-booking/internal/config"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWT(cfg *config.Config) func(http.Handler) http.Handler {
	// 🔑 Parse RSA public key ONCE
	keyPem := strings.ReplaceAll(cfg.JWTPublicKey, `\n`, "\n")

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyPem))
	if err != nil {
		panic("invalid JWT public key: " + err.Error())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return publicKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"]
			if !ok {
				http.Error(w, "Missing subject", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
