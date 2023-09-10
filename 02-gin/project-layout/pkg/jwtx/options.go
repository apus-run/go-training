package jwtx

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Option is config option.
type Option func(*Options) error

type Options struct {
	secretKey string

	ssid      string
	uid       uint64
	userAgent string

	// 集成 jwt.RegisteredClaims 配置
	jwt.RegisteredClaims
}

// DefaultOptions .
func DefaultOptions() *Options {
	now := time.Now()
	// 过期时间
	expiresAt := now.Add(time.Hour * 24)
	return &Options{
		secretKey: SecretKey,
		RegisteredClaims: jwt.RegisteredClaims{
			// 令牌过期时间
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			// 令牌的发行时间
			IssuedAt: jwt.NewNumericDate(now),
			// 令牌生效时间
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil
		}
	}
	return options
}

// WithSecretKey .
func WithSecretKey(secretKey string) Option {
	return func(o *Options) error {
		if secretKey == "" {
			return errors.New("secretKey can not be empty")
		}
		o.secretKey = secretKey
		return nil
	}
}

// WithUid .
func WithUid(uid uint64) Option {
	return func(o *Options) error {
		if uid == 0 {
			return errors.New("UID can not be empty")
		}
		o.uid = uid
		return nil
	}
}

// WithSsid .
func WithSsid(ssid string) Option {
	return func(o *Options) error {
		if len(ssid) == 0 {
			return errors.New("SSID can not be empty")
		}
		o.ssid = ssid
		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(o *Options) error {
		if ua == "" {
			return errors.New("UserAgent can not be empty")
		}
		o.userAgent = ua
		return nil
	}
}

// WithJwtRegisteredClaims  自行配置RegisteredClaims
func WithJwtRegisteredClaims(f func(options *Options)) Option {
	return func(config *Options) error {
		f(config)
		return nil
	}
}
