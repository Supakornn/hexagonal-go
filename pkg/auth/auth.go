package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/users"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type auth struct {
	mapClaims *mapClaims
	cfg       config.IJwtConfig
}

type admin struct {
	*auth
}

type apiKey struct {
	*auth
}

type mapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IAuth interface {
	SignToken() string
}

type IAdmin interface {
	SignToken() string
}

type IApiKey interface {
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

func (a *apiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.ApiKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method error:")
		}
		return cfg.SecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid: %v", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired: %v", err)
		} else {
			return nil, fmt.Errorf("token error %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token claims error")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method error:")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid: %v", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired: %v", err)
		} else {
			return nil, fmt.Errorf("token error %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token claims error")
	}
}

func ParseApiKey(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method error:")
		}
		return cfg.ApiKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid: %v", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired: %v", err)
		} else {
			return nil, fmt.Errorf("token error %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token claims error")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecom-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewAuth(TokenType TokenType, cfg config.IJwtConfig, Claims *users.UserClaims) (IAuth, error) {
	switch TokenType {
	case Access:
		return newAccessToken(cfg, Claims), nil
	case Refresh:
		return newRefreshToken(cfg, Claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("token type not found")
	}
}

func newAccessToken(cfg config.IJwtConfig, Claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: Claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecom-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, Claims *users.UserClaims) IAuth {
	return &auth{
		cfg: cfg,
		mapClaims: &mapClaims{
			Claims: Claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecom-api",
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
			mapClaims: &mapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ecom-api",
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

func newApiKey(cfg config.IJwtConfig) IAuth {
	return &apiKey{
		auth: &auth{
			cfg: cfg,
			mapClaims: &mapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ecom-api",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
