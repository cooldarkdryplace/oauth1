package etsy

import (
	"context"
	"net/http"
)

type userInfo struct {
	ID        int64  `json:"user_id"`
	LoginName string `json:"login_name"`
}

type userInfoResponse struct {
	Count   int        `json:"count"`
	Results []userInfo `json:"results"`
	Type    string     `json:"type"`
}

type EtsyClient struct {
	oauth *OAuth
}

func NewClient(o *OAuth) *EtsyClient {
	return &EtsyClient{oauth: o}
}

func (e *EtsyClient) Login(ctx context.Context) (string, TokenDetails, error) {
	creds, err := e.oauth.Login(ctx)
	if err != nil {
		return "", TokenDetails{}, err
	}

	details := TokenDetails{
		Token:       creds.OAuthToken,
		TokenSecret: creds.TokenSecret,
	}

	return creds.LoginURL, details, nil
}

func (e *EtsyClient) Callback(ctx context.Context, pin, token, secret string) (TokenDetails, error) {
	params := CallbackParams{
		Verifier:    pin,
		Token:       token,
		TokenSecret: secret,
	}

	return e.oauth.CompleteLogin(ctx, params)
}

func (e *EtsyClient) HTTPClient(accessToken, accessSecret string) *http.Client {
	return e.oauth.Client(accessToken, accessSecret)
}
