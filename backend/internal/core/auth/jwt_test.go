package core_auth

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestIssueAndParseToken(t *testing.T) {
	secret := "test-secret-key-at-least-32-bytes-long!"
	userID := 42
	email := "test@example.com"
	ttl := 15 * time.Minute

	token, err := IssueToken(secret, userID, email, ttl)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	// Token must have 3 dot-separated parts.
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %q, want %q", claims.Email, email)
	}
}

// TestClaimsAreNumericDate verifies that iat and exp are serialized as
// integer Unix timestamps (RFC 7519 NumericDate), not ISO-8601 strings.
func TestClaimsAreNumericDate(t *testing.T) {
	secret := "test-secret-key-at-least-32-bytes-long!"
	token, err := IssueToken(secret, 1, "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	parts := strings.Split(token, ".")
	claimsJSON, err := base64URLDecode(parts[1])
	if err != nil {
		t.Fatalf("decode claims: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(claimsJSON, &raw); err != nil {
		t.Fatalf("unmarshal raw claims: %v", err)
	}

	// iat and exp must be JSON numbers, not strings.
	for _, field := range []string{"iat", "exp"} {
		val := string(raw[field])
		if strings.HasPrefix(val, `"`) {
			t.Errorf("%s is a JSON string (%s), expected a JSON number (Unix timestamp)", field, val)
		}

		var n int64
		if err := json.Unmarshal(raw[field], &n); err != nil {
			t.Errorf("%s is not a valid JSON number: %v", field, err)
		}
	}
}

func TestParseTokenExpired(t *testing.T) {
	secret := "test-secret-key-at-least-32-bytes-long!"

	// Issue a token with a TTL of 0 (already expired).
	token, err := IssueToken(secret, 1, "a@b.com", -1*time.Second)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	_, err = ParseToken(secret, token)
	if err == nil {
		t.Fatal("ParseToken() expected error for expired token, got nil")
	}
}

func TestParseTokenInvalidSignature(t *testing.T) {
	token, err := IssueToken("secret-a", 1, "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	// Try parsing with a different secret.
	_, err = ParseToken("secret-b", token)
	if err == nil {
		t.Fatal("ParseToken() expected error for wrong secret, got nil")
	}
}

func TestParseTokenMalformed(t *testing.T) {
	_, err := ParseToken("secret", "not-a-jwt")
	if err == nil {
		t.Fatal("ParseToken() expected error for malformed token, got nil")
	}
}
