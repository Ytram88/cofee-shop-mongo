package auth

import (
	"cofee-shop-mongo/internal/utils"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"

	"time"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

var secret string

func SetSecret(secret string) {
	secret = secret
}

func WithJWTAuth(requiredRole []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get token from request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing token"))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := validateJWT(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		//extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			return
		}
		userID, ok := claims["sub"].(string)
		role, roleOk := claims["role"].(string)
		if !ok || !roleOk {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token data"))
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		accessFlag := false
		for _, r := range requiredRole {
			if role == r {
				accessFlag = true
			}
		}
		if !ok || !accessFlag {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("access denied: required role %s", requiredRole))
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func CreateJWT(userID, role string, expiration int64) (string, error) {
	expirationInSeconds := time.Second * time.Duration(expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(expirationInSeconds).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
