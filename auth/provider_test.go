package auth

import (
	"context"
	"testing"

	scyauth "github.com/viant/scy/auth"
	"github.com/viant/scy/auth/jwt"
	"github.com/viant/scy/auth/jwt/signer"
	"github.com/viant/scy/auth/jwt/verifier"
)

type credentialAuthenticator struct{}

func (credentialAuthenticator) BasicAuth(context.Context, string, string) (*scyauth.Token, error) {
	return nil, nil
}
func (credentialAuthenticator) VerifyIdentity(context.Context, string) (*jwt.Claims, error) {
	return nil, nil
}
func (credentialAuthenticator) ReissueIdentityToken(context.Context, string, string) (*scyauth.Token, error) {
	return nil, nil
}
func (credentialAuthenticator) ResetCredentials(context.Context, string, string) error { return nil }

type authProvider struct{ authenticator Authenticator }

func (p authProvider) Authenticator(Vendor) (Authenticator, error) { return p.authenticator, nil }
func (authProvider) Signer() *signer.Service                       { return nil }
func (authProvider) Verifier() *verifier.Service                   { return nil }

func TestProviderContractsCompile(t *testing.T) {
	var _ Authenticator = credentialAuthenticator{}
	var _ Auth = authProvider{authenticator: credentialAuthenticator{}}
	if ProviderKind != "auth" {
		t.Fatalf("ProviderKind = %q", ProviderKind)
	}
}
