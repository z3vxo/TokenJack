package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type Auth struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func New(accessSec, refreshSec []byte) *Auth {
	return &Auth{
		accessSecret:  accessSec,
		refreshSecret: refreshSec,
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
	}
}

func (a *Auth) Generate(username string) (Tokens, error) {
	access, err := a.sign(username, a.accessSecret, a.accessTTL)
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := a.sign(username, a.refreshSecret, a.refreshTTL)
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{AccessToken: access, RefreshToken: refresh}, nil

}

func (a *Auth) sign(username string, sec []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(sec)
}
