package main

import (
	"GoPrac/lesson2/crawler"
	"GoPrac/lesson2/process"
	"GoPrac/lesson2/requester"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const startUrl = "https://www.w3.org/Consortium/"

func main() {
	pid := os.Getpid()
	fmt.Printf("My PID is: %d\n", pid)

	//var r requester.Requester
	r := requester.NewRequester(time.Minute)
	//r = requester.NewRequestWithDelay(10*time.Second, r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	crawler := crawler.NewCrawler(2, r)
	crawler.Scan(ctx, startUrl, 0, false)

	chSig := make(chan os.Signal, 10)
	signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)

	go process.ProcessResult(ctx, crawler.GetResultChan(), cancel)

	go func() {
		for {
			sig := <-chSig
			switch sig {
			case syscall.SIGTERM:
				fmt.Printf("Signal SIGTERM caught\n")
				cancel()
			case syscall.SIGUSR1:
				crawler.MaxDepth += 2
			}
		}
	}()

	if <-ctx.Done() == struct{}{} {
		fmt.Printf("context canceled\n")
	}
}
