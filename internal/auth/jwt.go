package auth

import (
	"cofee-shop-mongo/internal/utils"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"

	"time"
)

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

}

func CreateJWT(secret []byte, userID string, expiration int64) (string, error) {
	expirationInSeconds := time.Second * time.Duration(expiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expirationInSeconds).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {

}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
}
