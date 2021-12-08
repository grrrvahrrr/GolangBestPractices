package page

import (
	"io"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

//Page package

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return page{}, err
	}
	return page{doc}, nil
}

func (p page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p page) GetLinks() []string {
	var urls []string
	startUrl := "https://www.w3.org"

	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		myUrl, ok := s.Attr("href")
		if ok {
			baseUrl, err := url.Parse(myUrl)
			if err != nil {
				log.Println("Couldn't parse myURL")
			}
			u, err := url.Parse("")
			if err != nil {
				log.Println("Couldn't parse myURL")
			}
			if baseUrl.IsAbs() {
				urls = append(urls, myUrl)
			} else {
				urls = append(urls, startUrl+baseUrl.ResolveReference(u).String())
			}
		}
	})
	return urls
}
