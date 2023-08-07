package ginx

import (
	"strings"

	"github.com/pkg/errors"
)

// Option is config option.
type Option func(*Options) error

type Options struct {
	mode         string // dev or prod
	addr         string
	port         string
	maxPingCount int
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		mode:         "dev",
		addr:         "localhost",
		port:         "8080",
		maxPingCount: 5,
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

// WithMode .
func WithMode(mode string) Option {
	return func(o *Options) error {
		if strings.ToLower(mode) != "dev" && strings.ToLower(mode) != "prod" {
			return errors.New("mode must be dev or prod")
		}
		o.mode = mode
		return nil
	}
}

// WithAddr .
func WithAddr(addr string) Option {
	return func(o *Options) error {
		if addr == "" {
			return errors.New("addr can not be empty")
		}
		o.addr = addr
		return nil
	}
}

// WithPort .
func WithPort(port string) Option {
	return func(o *Options) error {
		if port == "" {
			return errors.New("port can not be empty")
		}
		o.port = port
		return nil
	}
}

// WithMaxPingCount .
func WithMaxPingCount(maxPingCount int) Option {
	return func(o *Options) error {
		if maxPingCount <= 0 {
			return errors.New("maxPingCount must be greater than 0")
		}
		o.maxPingCount = maxPingCount
		return nil
	}
}
