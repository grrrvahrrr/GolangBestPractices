package requester

import (
	"context"
	"net/http"
	"time"
)

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type Request struct {
	Timeout time.Duration
}

func NewRequester(timeout time.Duration) Request {
	return Request{Timeout: timeout}
}

func (r Request) Get(ctx context.Context, url string) (Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		cl := &http.Client{
			Timeout: r.Timeout,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		body, err := cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer body.Body.Close()
		page, err := NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return page, nil
	}
}
