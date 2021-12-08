package requester

import (
	"GoPrac/lesson2/page"
	"context"
	"net/http"
	"time"
)

//Requester Package
type Requester interface {
	GetPage(ctx context.Context, url string) (page.Page, error)
}

// type reqWithDelay struct {
// 	delay time.Duration
// 	req   Requester
// }

// func NewRequestWithDelay(delay time.Duration, req Requester) *reqWithDelay {
// 	return &reqWithDelay{delay: delay, req: req}
// }

// func (r reqWithDelay) GetPage(ctx context.Context, url string) (page.Page, error) {
// 	time.Sleep(r.delay)
// 	return r.req.GetPage(ctx, url)
// }

/*
type HttpClient interface {
	 Do(r *http.Request) (*http.Response, error)
}
*/

type requester struct {
	timeout time.Duration
}

func NewRequester(timeout time.Duration) *requester {
	return &requester{timeout: timeout}
}

func (r requester) GetPage(ctx context.Context, url string) (page.Page, error) {
	cl := &http.Client{
		Timeout: r.timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	rawPage, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawPage.Body.Close()
	return page.NewPage(rawPage.Body)
}
