package auth

import (
	"context"
	"errors"
)

// Vendor identifies one supported authentication backend family.
type Vendor string

const (
	VendorJWT      Vendor = "jwt"
	VendorCognito  Vendor = "cognito"
	VendorFirebase Vendor = "firebase"
)

const (
	ClaimSubject = "sub"
	ClaimEmail   = "email"
	ClaimScope   = "scope"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// Claims carries normalized auth claims without tying the public surface to a
// specific vendor payload type.
type Claims map[string]interface{}

// Token is the extracted credential payload presented to an authenticator.
// Transport-specific decoration such as the Authorization header name or
// "Bearer " prefix should be removed before authentication.
type Token struct {
	Value string
}

// Principal is the reduced authenticated identity returned by an
// Authenticator.
type Principal struct {
	Subject string
	Email   string
	Claims  Claims
}

// Authenticator is the reduced stdlib-only authentication contract.
type Authenticator interface {
	Authenticate(ctx context.Context, token Token) (Principal, error)
}
