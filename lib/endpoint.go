package lib

import (
	"context"
	"net/http"
)

type Endpoint func(ctx context.Context, requestParam interface{}) (respenseResult interface{}, err error)

// EncodeRequestFunc 决定请求path方法
type EncodeRequestFunc func(ctx context.Context, r *http.Request, value interface{}) error
