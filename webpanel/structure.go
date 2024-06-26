package webpanel

import (
	"context"
	"encoding/xml"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"github.com/coreos/go-oidc/v3/oidc"
)

type WebPanel struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
	Salt []byte
	Config
	AuthConfig *AppAuthConfig
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type OIDCConfig struct {
	XMLName xml.Name `xml:"oidc"`
	ClientID string `xml:"clientID"`
	ClientSecret string `xml:"clientSecret"`
	RedirectURL string `xml:"redirectURL"`
	Scopes []string `xml:"scopes"`
	Provider string `xml:"provider"`
}

type Config struct {
	Username        string `xml:"username"`
	Password        string `xml:"password"`
	DatabaseAddress string `xml:"databaseAddress"`
	DatabaseName    string `xml:"databaseName"`
	Address         string `xml:"address"`
	AssetsPath      string `xml:"assetsPath"`
	OIDCConfig	    OIDCConfig `xml:"oidc"`
	AuthMode 	  	bool `xml:"auth_mode"`
}

type AppAuthConfig struct {
    OAuth2Config *oauth2.Config
    Provider     *oidc.Provider
}
