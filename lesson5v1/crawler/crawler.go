package crawler

import (
	"context"
	"sync"
)

type CrawlResult struct {
	Err   error
	Title string
	Url   string
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, depth int)
	ChanResult() <-chan CrawlResult
}

type Crawl struct {
	R       Requester
	Res     chan CrawlResult
	Visited map[string]struct{}
	Mu      sync.RWMutex
}

func NewCrawler(r Requester) *Crawl {
	return &Crawl{
		R:       r,
		Res:     make(chan CrawlResult),
		Visited: make(map[string]struct{}),
		Mu:      sync.RWMutex{},
	}
}

func (c *Crawl) ChanResult() <-chan CrawlResult {
	return c.Res
}

func (c *Crawl) Scan(ctx context.Context, url string, depth int) {
	if depth <= 0 { //Проверяем то, что есть запас по глубине
		return
	}
	c.Mu.RLock()
	_, ok := c.Visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.Mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.R.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.Res <- CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.Mu.Lock()
		c.Visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.Mu.Unlock()
		c.Res <- CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth-1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}
