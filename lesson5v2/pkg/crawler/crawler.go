package crawler

import (
	"GoBest/lesson5v2/pkg/domain"
	"context"
	"sync"
)

type crawler struct {
	r       domain.Requester
	res     chan domain.CrawlResult
	visited map[string]struct{}
	mu      sync.RWMutex
}

func NewCrawler(r domain.Requester) *crawler {
	return &crawler{
		r:       r,
		res:     make(chan domain.CrawlResult),
		visited: make(map[string]struct{}),
		mu:      sync.RWMutex{},
	}
}

func (c *crawler) Scan(ctx context.Context, url string, depth int) {
	if depth <= 0 { //Проверяем то, что есть запас по глубине
		return
	}
	c.mu.RLock()
	_, ok := c.visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.res <- domain.CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.mu.Lock()
		c.visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.mu.Unlock()
		c.res <- domain.CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}
		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, depth-1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

func (c *crawler) ChanResult() <-chan domain.CrawlResult {
	return c.res
}
