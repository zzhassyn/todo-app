// Package core_auth provides a minimal, dependency-free JSON Web Token
// implementation (HS256 only) used to authenticate HTTP requests across
// features. It deliberately avoids pulling in a third-party JWT library:
// the token shape needed by this application is small enough that a
// hand-rolled, well-tested implementation keeps the dependency graph
// minimal while remaining fully interoperable with the JWT spec (a token
// minted here can be decoded by any standard JWT library, and vice versa,
// as long as the algorithm is HS256).
package core_auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

const tokenAlgorithm = "HS256"
const tokenType = "JWT"

type header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// Claims represents a JWT payload. IssuedAt and ExpiresAt are stored as
// Unix epoch seconds (NumericDate per RFC 7519 §2) so that tokens are
// interoperable with any standard JWT library.
type Claims struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func base64URLEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// IssueToken creates a signed JWT for the given user, valid for ttl.
func IssueToken(secret string, userID int, email string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	headerJSON, err := json.Marshal(header{Algorithm: tokenAlgorithm, Type: tokenType})
	if err != nil {
		return "", fmt.Errorf("marshal jwt header: %w", err)
	}

	claims := Claims{
		UserID:    userID,
		Email:     email,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal jwt claims: %w", err)
	}

	signingInput := base64URLEncode(headerJSON) + "." + base64URLEncode(claimsJSON)
	signature := sign(secret, signingInput)

	return signingInput + "." + base64URLEncode(signature), nil
}

// ParseToken validates the token's signature and expiry, returning its claims.
func ParseToken(secret string, token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("malformed token: expected 3 parts, got %d: %w", len(parts), core_errors.ErrInvalidArgument)
	}

	signingInput := parts[0] + "." + parts[1]

	gotSignature, err := base64URLDecode(parts[2])
	if err != nil {
		return Claims{}, fmt.Errorf("decode signature: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	wantSignature := sign(secret, signingInput)
	if subtle.ConstantTimeCompare(gotSignature, wantSignature) != 1 {
		return Claims{}, fmt.Errorf("invalid token signature: %w", core_errors.ErrInvalidArgument)
	}

	claimsJSON, err := base64URLDecode(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("decode claims: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return Claims{}, fmt.Errorf("unmarshal claims: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	if time.Now().UTC().Unix() > claims.ExpiresAt {
		return Claims{}, fmt.Errorf("token expired at %d: %w", claims.ExpiresAt, core_errors.ErrInvalidArgument)
	}

	return claims, nil
}

func sign(secret string, signingInput string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	return mac.Sum(nil)
}

