package main

import (
	"GoBest/lesson5v2/pkg/config"
	"GoBest/lesson5v2/pkg/crawler"
	"GoBest/lesson5v2/pkg/processor"
	"GoBest/lesson5v2/pkg/requester"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	run(ctx, cancel)

	sigCh := make(chan os.Signal, 10)    //Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT) //Подписываемся на сигнал SIGINT
	for {
		select {
		case <-ctx.Done(): //Если всё завершили - выходим
			return
		case <-sigCh:
			cancel() //Если пришёл сигнал SigInt - завершаем контекст
		}
	}
}

func run(ctx context.Context, cancel func()) {
	cfg := config.Config{
		MaxDepth:   2,
		MaxResults: 10,
		MaxErrors:  5,
		Url:        "https://www.glassleaf.ru/",
		Timeout:    10,
	}

	r := requester.NewRequester(time.Duration(cfg.Timeout) * time.Second)
	cr := crawler.NewCrawler(r)

	go cr.Scan(ctx, cfg.Url, cfg.MaxDepth)           //Запускаем краулер в отдельной рутине
	go processor.ProcessResult(ctx, cancel, cr, cfg) //Обрабатываем результаты в отдельной рутине
}
