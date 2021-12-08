package crawler

import (
	"GoPrac/lesson2/requester"
	"context"
	"sync"
)

//Crawler package

type Crawler interface {
	Scan(ctx context.Context, url string, curDepth int)
	GetResultChan() <-chan CrawlResult
}

type CrawlResult struct {
	Title string
	Url   string
	Err   error
}

type Crawl struct {
	MaxDepth  int
	infDepth  bool
	req       requester.Requester
	res       chan CrawlResult
	visited   map[string]struct{}
	visitedMu sync.RWMutex
}

func (c *Crawl) GetResultChan() <-chan CrawlResult {
	return c.res
}

func NewCrawler(maxDepth int, req requester.Requester) *Crawl {
	return &Crawl{
		MaxDepth: maxDepth,
		req:      req,
		res:      make(chan CrawlResult, 100),
		visited:  make(map[string]struct{}),
	}
}

func (c *Crawl) Scan(ctx context.Context, url string, curDepth int, infDepth bool) {
	c.infDepth = infDepth
	c.visitedMu.RLock()
	if _, ok := c.visited[url]; ok {
		c.visitedMu.RUnlock()
		return
	}
	c.visitedMu.RUnlock()
	if curDepth >= c.MaxDepth && !c.infDepth {
		return
	}
	select {
	case <-ctx.Done():
		return
	default:
		page, err := c.req.GetPage(ctx, url)
		c.visitedMu.Lock()
		c.visited[url] = struct{}{}
		c.visitedMu.Unlock()
		if err != nil {
			c.res <- CrawlResult{Url: url, Err: err}
			return
		}
		title := page.GetTitle()
		c.res <- CrawlResult{
			Title: title,
			Url:   url,
			Err:   nil,
		}
		links := page.GetLinks()
		for _, link := range links {
			go c.Scan(ctx, link, curDepth+1, c.infDepth)
		}
	}
}
