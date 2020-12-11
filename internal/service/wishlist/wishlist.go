package wishlist

import (
	"context"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"github.com/satori/go.uuid"
	"time"
)

type WishlistRepoIFace interface {
	Insert(wishlist Wish) error
	Fetch(wishID string) (*Wish, error)
	Update(wishID string, originalPrice, currentPrice int64, discountRate float64) error
	Delete(wishID string) error
	FetchByCustomer(customerRefID string) ([]Wish, error)
	FetchAll() ([]Wish, error)
}

type Wish struct {
	WishID        string     `json:"id"`
	CustomerRefID string     `json:"customer_ref_id"`
	ProductName   string     `json:"product_name"`
	CurrentPrice  int64      `json:"current_price"`
	OriginalPrice int64      `json:"original_price"`
	DiscountRate  float64    `json:"discount_rate"`
	Source        string     `json:"source"`
	UpdatedAt     *time.Time `json:"last_seen_at"`
}

type wishlistSvc struct {
	repo          WishlistRepoIFace
	productClient product.Client
}

func NewWishlistSvc(repo WishlistRepoIFace, productClient product.Client) *wishlistSvc {
	return &wishlistSvc{
		repo:          repo,
		productClient: productClient,
	}
}

func (w *wishlistSvc) Add(wishlist *Wish) error {
	item, err := w.productClient.GetItem(context.Background(), wishlist.Source)
	if err != nil {
		return err
	}

	id := uuid.NewV4()

	wishlist.WishID = id.String()
	wishlist.ProductName = item.Name
	wishlist.OriginalPrice = item.OriginalPrice
	wishlist.CurrentPrice = item.CurrentPrice
	wishlist.DiscountRate = item.DiscountRate

	err = w.repo.Insert(*wishlist)
	if err != nil {
		return err
	}
	return nil
}

func (w *wishlistSvc) FetchByCustomer(customerRefID string) ([]Wish, error) {
	wishlist, err := w.repo.FetchByCustomer(customerRefID)
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (w *wishlistSvc) FetchAll() ([]Wish, error) {
	wishlist, err := w.repo.FetchAll()
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (w *wishlistSvc) DeleteWish(wishID string) error {
	err := w.repo.Delete(wishID)
	if err != nil {
		return err
	}
	return nil
}
