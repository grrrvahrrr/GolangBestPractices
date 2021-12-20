package main

import (
	"GoBest/lesson5v1/config"
	"GoBest/lesson5v1/connector"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//MAIN
func main() {

	cfg := config.NewConfig(2, 10, 5, "https://www.glassleaf.ru/", 10)

	r := connector.NewRequester(time.Duration(cfg.Timeout) * time.Second)
	cr := connector.NewCrawler(r)

	ctx, cancel := context.WithCancel(context.Background())
	go cr.Scan(ctx, cfg.Url, cfg.MaxDepth)            //Запускаем краулер в отдельной рутине
	go connector.ProcessResult(ctx, cancel, cr, *cfg) //Обрабатываем результаты в отдельной рутине

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
