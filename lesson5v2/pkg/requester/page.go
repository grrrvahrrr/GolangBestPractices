package requester

import (
	"GoBest/lesson5v2/pkg/domain"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (domain.Page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return nil, err
	}
	return &page{doc: doc}, nil
}

func (p *page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
		}
	})
	return urls
}
