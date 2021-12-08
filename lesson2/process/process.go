package process

import (
	"GoPrac/lesson2/crawler"
	"context"
	"fmt"
)

//Process package

func ProcessResult(ctx context.Context, in <-chan crawler.CrawlResult, cancel context.CancelFunc) {
	var errCount int
	for {
		select {
		case res := <-in:
			if res.Err != nil {
				errCount++
				fmt.Printf("ERROR Link: %s, err: %v\n", res.Url, res.Err)
				if errCount >= 3 {
					cancel()
				}
			} else {
				fmt.Printf("Link: %s, Title: %s\n", res.Url, res.Title)
			}
		case <-ctx.Done():
			fmt.Printf("context canceled in process.go\n")
			return
		}
	}
}
