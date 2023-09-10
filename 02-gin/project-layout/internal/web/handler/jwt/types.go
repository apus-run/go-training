package jwt

import (
	"project-layout/pkg/ginx"
	"project-layout/pkg/jwtx"
)

type Handler interface {
	SetLoginToken(ctx *ginx.Context, uid uint64) error
}

type UserClaims struct {
	jwtx.CustomClaims
}

type RefreshClaims struct {
	jwtx.CustomClaims
}
