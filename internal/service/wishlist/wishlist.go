package wishlist

import (
	"context"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"github.com/satori/go.uuid"
)

type wishlistRepoIFace interface {
	Insert(wishlist Wishlist) error
	FetchByCustomer(customerRefID string) ([]Wishlist, error)
	FetchAll() ([]Wishlist, error)
}

type Wishlist struct {
	WishlistID    string `json:"id"`
	CustomerRefID string `json:"customer_ref_id"`
	ProductName   string `json:"product_name"`
	CurrentPrice  int64  `json:"current_price"`
	OriginalPrice int64  `json:"original_price"`
	Source        string `json:"source"`
}

type FetchByCustomerResp struct {
	Wishlists []Wishlist `json:"data"`
}

type wishlistSvc struct {
	repo          wishlistRepoIFace
	productClient product.Client
}

func NewWishlistSvc(repo wishlistRepoIFace, productClient product.Client) *wishlistSvc {
	return &wishlistSvc{
		repo:          repo,
		productClient: productClient,
	}
}

func (w *wishlistSvc) Add(wishlist *Wishlist) error {
	item, err := w.productClient.GetItem(context.Background(), wishlist.Source)
	if err != nil {
		return err
	}

	id := uuid.NewV4()

	wishlist.WishlistID = id.String()
	wishlist.ProductName = item.Name
	wishlist.OriginalPrice = item.OriginalPrice
	wishlist.CurrentPrice = item.CurrentPrice

	err = w.repo.Insert(*wishlist)
	if err != nil {
		return err
	}
	return nil
}

func (w *wishlistSvc) FetchByCustomer(customerRefID string) ([]Wishlist, error) {
	wishlist, err := w.repo.FetchByCustomer(customerRefID)
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (w *wishlistSvc) FetchAll() ([]Wishlist, error) {
	wishlist, err := w.repo.FetchAll()
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}
