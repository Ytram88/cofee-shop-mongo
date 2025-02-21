package jtoken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	expirationJWT = time.Hour * 5
	HeaderJWT     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	/*
		{
			"alg": "HS256",
			"typ": "JWT"
		}
	*/
)

type Payload struct {
	IsAdmin   bool   `json:"is_admin"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewPayload() *Payload {
	return &Payload{}
}

// VerifyJWT checks the token format, signature, and expiration
func VerifyJWT(token, secretKey string) (*Payload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	if parts[0] != HeaderJWT {
		return nil, errors.New("invalid header format")
	}

	// Decode payload from Base64
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	payload := NewPayload()
	if err := json.Unmarshal(payloadJSON, payload); err != nil {
		return nil, errors.New("invalid payload format")
	}

	// Check expiration
	if payload.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}

	// Verify Signature
	expectedSignature, err := SignHS256(parts[0]+"."+parts[1], secretKey)
	if err != nil {
		return nil, err
	}

	if parts[2] != expectedSignature {
		return nil, errors.New("invalid signature")
	}

	return payload, nil
}

// SignHS256 generates an HMAC-SHA256 signature
func SignHS256(data string, secretKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := h.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(signature), nil
}

// GenerateAccessToken creates a signed JWT
func GenerateAccessToken(payload *Payload, secretKey string) (string, error) {
	payload.ExpiresAt = time.Now().Add(expirationJWT).Unix()

	data, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("invalid payload")
	}

	// Encode payload in Base64
	encodedPayload := base64.RawURLEncoding.EncodeToString(data)

	// Create signature
	signature, err := SignHS256(HeaderJWT+"."+encodedPayload, secretKey)
	if err != nil {
		return "", fmt.Errorf("unable to generate access token: %v", err)
	}

	jwtToken := HeaderJWT + "." + encodedPayload + "." + signature

	return jwtToken, nil
}
