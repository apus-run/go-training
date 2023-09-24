package ratelimit

import (
	"context"
	"project-layout/internal/service/sms"
	"project-layout/pkg/ratelimit_redis"
	"testing"
)

func TestRatelimitSMSService_Send(t *testing.T) {
	type fields struct {
		svc     sms.Service
		limiter ratelimit_redis.Limiter
	}
	type args struct {
		ctx     context.Context
		tplId   string
		args    []string
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RatelimitSMSService{
				svc:     tt.fields.svc,
				limiter: tt.fields.limiter,
			}
			if err := r.Send(tt.args.ctx, tt.args.tplId, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
