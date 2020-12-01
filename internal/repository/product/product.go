package product

import "context"

type Client interface {
	GetItem(ctx context.Context, source string) (item *Item, err error)
}

type Item struct {
	Name          string
	CurrentPrice  int64
	OriginalPrice int64
	Source        string
}
