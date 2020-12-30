package ctxlib

import (
	"context"
)

type key string

var (
	privilegedKey key = "previleged-key"
)

func SetPrevilegedAccessKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, privilegedKey, key)
}

func PrevilegedAccessKey(ctx context.Context) string {
	val, ok := ctx.Value(privilegedKey).(string)
	if !ok {
		return ""
	}
	return val
}
