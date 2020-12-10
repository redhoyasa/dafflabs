package pricealert

import (
	"context"
	"fmt"
	qm "github.com/quickmetrics/qckm-go"
	"github.com/redhoyasa/dafflabs/internal/client/telegram"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
	"github.com/rs/zerolog/log"
)

type Client struct {
	telegram     telegram.Client
	productRepo  product.Client
	wishlistRepo wishlist.WishlistRepoIFace
}

func NewClient(telegram telegram.Client, productRepo product.Client, wishlistRepo wishlist.WishlistRepoIFace) (*Client, error) {
	c := new(Client)
	c.productRepo = productRepo
	c.telegram = telegram
	c.wishlistRepo = wishlistRepo
	return c, nil
}

func (c *Client) GeneratePriceChecker() error {
	wishes, err := c.wishlistRepo.FetchAll()
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	for _, wish := range wishes {
		go c.CheckPrice(context.Background(), wish.WishID)
	}

	return nil
}

// CheckPrice does:
// 1. Get current wish
// 2. Check product price
// 3. When new current price is lower than existing current price, send notification
// 4. Update wish
func (c *Client) CheckPrice(ctx context.Context, wishID string) error {
	log.Info().
		Str("WishID", wishID).
		Msg("Start checking price")

	wish, err := c.wishlistRepo.Fetch(wishID)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	item, err := c.productRepo.GetItem(ctx, wish.Source)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	if item.CurrentPrice < wish.CurrentPrice {
		msg := fmt.Sprintf("Product %s di %s harganya cuma %d", item.Name, item.Source, item.CurrentPrice)
		err = c.telegram.SendMessage(msg)

		if err != nil {
			log.Err(err).Msg(err.Error())
			return err
		}
	}

	err = c.wishlistRepo.Update(wishID, item.OriginalPrice, item.CurrentPrice, 0)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	log.Info().
		Str("WishID", wishID).
		Msg("Finish checking price")

	qm.Event("wishlist.checked", 1)

	return nil
}
