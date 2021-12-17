package connector

import (
	"GoBest/lesson5v1/crawler"
	"GoBest/lesson5v1/requester"
	"context"
	"time"
)

type Requester interface {
	Get(ctx context.Context, url string) (requester.Page, error)
}

func NewRequester(timeout time.Duration) Request {
	return Request{Timeout: timeout}
}

type Request struct {
	Timeout time.Duration
}

func (r Request) Get(ctx context.Context, url string) (crawler.Page, error) {
	var rr requester.Request
	return rr.Get(ctx, url)
}
