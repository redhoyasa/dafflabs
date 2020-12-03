package wishlist

type wishlistRepoIFace interface {
	Insert(wishlist Wishlist) error
	FetchByCustomer(customerRefID string) ([]Wishlist, error)
	FetchAll() ([]Wishlist, error)
}

type wishlistSvcIFace interface {
	Add(customerRefID string, source string) error
	List(customerRefID string) ([]Wishlist, error)
	Remove(customerRefID string, source string) error
}

type Wishlist struct {
	WishlistID    int64
	CustomerRefID string
	ProductName   string
	CurrentPrice  int64
	OriginalPrice int64
	Source        string
}

type wishlistSvc struct {
	repo wishlistRepoIFace
}
