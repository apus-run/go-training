package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"project-layout/internal/service/sms"
	"project-layout/pkg/jwtx"
)

type Service struct {
	svc sms.Service

	key string
}

// GenerateToken 手动生成一个token
func (s *Service) GenerateToken(ctx context.Context, tplId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, SMSClaims{
		TplId: tplId,
	})
	tokenString, err := token.SignedString([]byte(s.key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) Send(ctx context.Context,
	tplToken string, args []string, numbers ...string) error {
	var c SMSClaims

	// 如果解析(校验token)成功, 说明调用的业务方是合法的, 是我颁发的token
	token, err := jwtx.DecodeJWT(tplToken, s.key, &SMSClaims{})
	if err != nil {
		return err
	}
	// 是否合法
	if !token.Valid {
		return errors.New("token 不合法")
	}

	return s.svc.Send(ctx, c.TplId, args, numbers...)
}

type SMSClaims struct {
	TplId string

	jwt.RegisteredClaims
}
