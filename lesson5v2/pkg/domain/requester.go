package domain

import "context"

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}
