package handler

import (
	"context"

	jwt "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	tokenKey = contextKey("jwt")
)

// withToken adds the JWT token to the context.
func withToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// getToken reads the JWT token from the context.
func getToken(ctx context.Context) *jwt.Token {
	obj := ctx.Value(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func getClaims(ctx context.Context) *AccessTokenClaims {
	token := getToken(ctx)
	if token == nil {
		return nil
	}
	return token.Claims.(*AccessTokenClaims)
}
