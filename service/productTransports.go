package service

import (
	"net/http"
	"strconv"
	"context"
)

func ProdEncodeFunc(ctx context.Context, httpRequest *http.Request, requestParams interface{}) error {
	// 类型断言。
	prodr := requestParams.(ProductRequest)
	httpRequest.URL.Path+="/product/"+strconv.Itoa(prodr.ProductId)
	return nil
}
