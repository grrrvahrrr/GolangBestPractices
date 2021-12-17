package connector

import (
	"GoBest/lesson5v1/crawler"
	"context"
	"sync"
)

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, depth int)
	ChanResult() <-chan crawler.CrawlResult
}

func NewCrawler(r crawler.Requester) *crawler.Crawl {
	return &crawler.Crawl{
		R:       r,
		Res:     make(chan crawler.CrawlResult),
		Visited: make(map[string]struct{}),
		Mu:      sync.RWMutex{},
	}
}
