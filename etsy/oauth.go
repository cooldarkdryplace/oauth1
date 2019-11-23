package etsy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cooldarkdryplace/oauth1"
)

const (
	ListingsRead = "listings_r" // Read a members's inactive and expired (i.e., non-public) listings.
	ProfileRead  = "profile_r"  // Read a member's private profile information.
)

type TokenDetails struct {
	Token       string
	TokenSecret string
}

type OAuth struct {
	config oauth1.Config
}

func NewOAuth(ck, ss string) *OAuth {
	return &OAuth{
		config: oauth1.Config{
			ConsumerKey:    ck,
			ConsumerSecret: ss,
			Endpoint: oauth1.Endpoint{
				RequestTokenURL: "https://openapi.etsy.com/v2/oauth/request_token?scope=" + ListingsRead,
				AuthorizeURL:    "https://openapi.etsy.com/v2/oauth/access_token",
				AccessTokenURL:  "https://openapi.etsy.com/v2/oauth/access_token",
			},
		},
	}
}

func (oa *OAuth) Login(ctx context.Context) (oauth1.Credentials, error) {
	creds, err := oa.config.RequestToken()
	if err != nil {
		return oauth1.Credentials{}, fmt.Errorf("failed to get temp creds: %s", err)
	}

	return creds, nil
}

type CallbackParams struct {
	Verifier    string
	Token       string
	TokenSecret string
}

func (oa *OAuth) CompleteLogin(ctx context.Context, params CallbackParams) (TokenDetails, error) {
	accessToken, accessSecret, err := oa.config.AccessToken(params.Token, params.TokenSecret, params.Verifier)
	if err != nil {
		return TokenDetails{}, fmt.Errorf("failed to get AccessToken: %s", err)
	}

	details := TokenDetails{
		Token:       accessToken,
		TokenSecret: accessSecret,
	}

	return details, nil
}

func (oa *OAuth) Client(accessToken, accessSecret string) *http.Client {
	token := oauth1.NewToken(accessToken, accessSecret)
	return oa.config.Client(oauth1.NoContext, token)
}
