package core_auth

import "context"

type contextKey struct{}

var key = contextKey{}

// ToContext stores the authenticated user's claims in the context. Called
// by the Auth HTTP middleware after successfully validating the request's
// token.
func ToContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

// FromContext retrieves the authenticated user's claims from the context.
// The second return value is false if no authenticated user is present
// (i.e. the Auth middleware was not applied to this route, or the request
// was unauthenticated).
func FromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(key).(Claims)
	return claims, ok
}
