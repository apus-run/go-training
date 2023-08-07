package jwt

import (
	"github.com/pkg/errors"
	"time"
)

// Option is config option.
type Option func(*Options) error

type Options struct {
	secretKey string
	userId    uint64
	expireAt  time.Time
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		expireAt: time.Now(),
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

// WithUserId .
func WithUserId(userId uint64) Option {
	return func(o *Options) error {
		if userId == 0 {
			return errors.New("userId can not be empty")
		}
		o.userId = userId
		return nil
	}
}

// WithExpireAt .
func WithExpireAt(expireAt time.Time) Option {
	return func(o *Options) error {
		if expireAt.IsZero() {
			return errors.New("expireAt can not be empty")
		}
		o.expireAt = expireAt
		return nil
	}
}
