// context 的键值的存储并不是覆盖, 而是逐级存储
package main

import (
	"context"
)

type ctxKey struct{}

// NewContext returns a new Context that carries context.Context value.
func NewContext(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, ctxKey{}, value)
}

// FromContext returns the context.Context value stored in ctx, if any.
func FromContext(ctx context.Context) any {
	value := ctx.Value(ctxKey{})
	return value
}
