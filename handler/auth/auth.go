package auth

import (
	"context"
	"github.com/viant/scy/auth"
	"github.com/viant/scy/auth/jwt"
	"github.com/viant/scy/auth/jwt/signer"
	"github.com/viant/scy/auth/jwt/verifier"
)

// Vendor represents an authentication service vendor
type Vendor string

const (
	VendorCognito  Vendor = "cognito"
	VendorFirebase Vendor = "firebase"
)

// Auth represents an authentication service
type Auth interface {
	Authenticator(vendor Vendor) (Authenticator, error)
	Signer() *signer.Service
	Verifier() *verifier.Service
}

type Authenticator interface {
	BasicAuth(ctx context.Context, user string, password string) (*auth.Token, error)
	VerifyIdentity(ctx context.Context, idToken string) (*jwt.Claims, error)
	ReissueIdentityToken(ctx context.Context, refreshToken string, subject string) (*auth.Token, error)
}
