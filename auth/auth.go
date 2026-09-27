package auth

import (
	"context"
	"errors"

	scyauth "github.com/viant/scy/auth"
	"github.com/viant/scy/auth/jwt"
	"github.com/viant/scy/auth/jwt/signer"
	"github.com/viant/scy/auth/jwt/verifier"
)

// ProviderKind is the static binding key used by Datly custom handlers.
const ProviderKind = "auth"

// Vendor identifies one supported authentication backend family.
type Vendor string

const (
	VendorDefault  Vendor = "default"
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
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("expired token")
	ErrUnsupportedVendor = errors.New("unsupported auth vendor")
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

// TokenAuthenticator is the reduced token-to-principal contract used by
// generic protocol adapters.
type TokenAuthenticator interface {
	Authenticate(ctx context.Context, token Token) (Principal, error)
}

// Auth is the statically bound authentication provider used by application
// handlers. Implementations are composed by the host; xDatly owns only this
// contract and never loads providers dynamically.
type Auth interface {
	Authenticator(vendor Vendor) (Authenticator, error)
	Signer() *signer.Service
	Verifier() *verifier.Service
}

// Authenticator preserves the original Datly/xDatly credential-provider
// contract used by Auth's OAuth handlers.
type Authenticator interface {
	BasicAuth(ctx context.Context, user, password string) (*scyauth.Token, error)
	VerifyIdentity(ctx context.Context, idToken string) (*jwt.Claims, error)
	ReissueIdentityToken(ctx context.Context, refreshToken, subject string) (*scyauth.Token, error)
	ResetCredentials(ctx context.Context, email, newPassword string) error
}
