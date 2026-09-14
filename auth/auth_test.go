package auth

import (
	"context"
	"errors"
	"testing"
)

type testAuthenticator struct{}

func (testAuthenticator) Authenticate(_ context.Context, token Token) (Principal, error) {
	return Principal{
		Subject: token.Value,
		Claims:  Claims{ClaimSubject: token.Value},
	}, nil
}

func TestAuthContractCompile(t *testing.T) {
	var _ Authenticator = testAuthenticator{}
	if VendorJWT != "jwt" || VendorCognito != "cognito" || VendorFirebase != "firebase" {
		t.Fatalf("expected vendor constants to be preserved")
	}
	if !errors.Is(ErrInvalidToken, ErrInvalidToken) || !errors.Is(ErrExpiredToken, ErrExpiredToken) {
		t.Fatalf("expected exported auth sentinel errors")
	}
}

func TestAuthenticatorAuthenticate(t *testing.T) {
	authn := testAuthenticator{}
	principal, err := authn.Authenticate(context.Background(), Token{Value: "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if principal.Subject != "abc" {
		t.Fatalf("expected subject to be preserved, got %q", principal.Subject)
	}
	if got := principal.Claims[ClaimSubject]; got != "abc" {
		t.Fatalf("expected canonical subject claim, got %v", got)
	}
}
