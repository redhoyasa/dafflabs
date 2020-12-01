package pricealert

import (
	"context"
	"fmt"
	qm "github.com/quickmetrics/qckm-go"
	"github.com/redhoyasa/dafflabs/internal/client/telegram"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"github.com/rs/zerolog/log"
)

type Client struct {
	telegram    telegram.Client
	productRepo product.Client
}

func NewClient(telegram telegram.Client, productRepo product.Client) (*Client, error) {
	c := new(Client)
	c.productRepo = productRepo
	c.telegram = telegram
	return c, nil
}

func (c *Client) CheckPrice(ctx context.Context, productUrl string, threshold int64) {
	item, err := c.productRepo.GetItem(ctx, productUrl)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return
	}

	if item.OriginalPrice > item.CurrentPrice {
		msg := fmt.Sprintf("Product %s di %s harganya cuma %d", item.Name, item.Source, item.CurrentPrice)
		err = c.telegram.SendMessage(msg)

		if err != nil {
			log.Err(err).Msg(err.Error())
			return
		}
	}

	log.Info().
		Str("ProductUrl", productUrl).
		Msg("Finish checking price")

	qm.Event("wishlist.checked", 1)
}
