package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	Apikey  TokenType = "apikey"
)

type auth struct {
	mapClaims *authMapClaims
	cfg       config.IJwtConfig
}

type admin struct {
	*auth
}

type authMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IAuth interface {
	SignToken() string
}

type IAdmin interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *admin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("token is invalid")
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("token is invalid")
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "hexagonal-go",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case Apikey:
	default:
		return nil, fmt.Errorf("token type is invalid")
	}

	return nil, nil
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "hexagonal-go",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "hexagonal-go",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuth {
	return &admin{
		auth: &auth{
			cfg: cfg,
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "hexagonal-go",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
