package main

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	t.Logf("%s\n", ctx)
	t.Logf("Done \t%#v\n", ctx.Done())
	t.Logf("Err \t%#v\n", ctx.Err())
	t.Logf("Value \t%#v\n", ctx.Value("key"))
	deadline, ok := ctx.Deadline()
	t.Logf("Deadline \t%s(%t)\n", deadline, ok)

	t.Log("-----------------------------------------")

	// key 推荐使用 Type 定义, String 类型不安全
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key1", "value2")
	ctx = context.WithValue(ctx, "key1", "value3")
	t.Logf("Context \t%#v\n", ctx)
	t.Logf("Value \t%#v\n", ctx.Value("key1"))
}

const RequestID = "requestID"

func TestContextType(t *testing.T) {
	ctx := context.Background()
	SendContext(ctx)

}

type CtxSendKey string

func SendContext(ctx context.Context) {
	key := CtxSendKey(RequestID)
	ctx = context.WithValue(ctx, key, "123")
	ReciverContext(ctx)
}

type CtxReciverKey string

func ReciverContext(ctx context.Context) {
	key := CtxReciverKey(RequestID)
	ctx = context.WithValue(ctx, key, "456")
	LoggerContext(ctx)
}

func LoggerContext(ctx context.Context) {
	key1 := CtxSendKey(RequestID)
	key2 := CtxReciverKey(RequestID)
	fmt.Println("Send", key1, ctx.Value(key1))
	fmt.Println("Receiver", key2, ctx.Value(key2))
}
