package repository

import (
	"context"
	"fmt"
	"github.com/redhoyasa/dafflabs/internal/database"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
)

type WishlistRepo struct {
	db database.Database
}

func NewWishlistRepo(db database.Database) *WishlistRepo {
	return &WishlistRepo{
		db: db,
	}
}

func (w *WishlistRepo) Insert(wishlist wishlist.Wishlist) (err error) {
	query := fmt.Sprintf(
		`INSERT INTO wishlists(
					id,
					customer_ref_id, 
					product_name, 
					current_price, 
					original_price, 
					source
				) VALUES (
					$1, $2, $3, $4, $5, $6
				)`)

	tx, err := w.db.Begin()
	if err != nil {
		return
	}

	if _, err = tx.Exec(query, wishlist.WishlistID, wishlist.CustomerRefID, wishlist.ProductName, wishlist.CurrentPrice, wishlist.OriginalPrice, wishlist.Source); err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()
	return
}

func (w *WishlistRepo) FetchByCustomer(customerRefID string) (wishlists []wishlist.Wishlist, err error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			customer_ref_id,
			product_name,
			current_price,
			original_price,
			source
		FROM wishlists
		WHERE
			is_deleted = 'false'
			AND customer_ref_id = $1
	`)

	rows, err := w.db.QueryContext(context.Background(), query, customerRefID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		wishlist := wishlist.Wishlist{}
		err := rows.Scan(
			&wishlist.WishlistID,
			&wishlist.CustomerRefID,
			&wishlist.ProductName,
			&wishlist.CurrentPrice,
			&wishlist.OriginalPrice,
			&wishlist.Source)

		if err != nil {
			return nil, err
		}

		wishlists = append(wishlists, wishlist)
	}
	return
}

func (w *WishlistRepo) FetchAll() (wishlists []wishlist.Wishlist, err error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			customer_ref_id,
			product_name,
			current_price,
			original_price,
			source
		FROM wishlists
		WHERE
			is_deleted = 'false'
	`)

	rows, err := w.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		wishlist := wishlist.Wishlist{}
		err := rows.Scan(
			&wishlist.WishlistID,
			&wishlist.CustomerRefID,
			&wishlist.ProductName,
			&wishlist.CurrentPrice,
			&wishlist.OriginalPrice,
			&wishlist.Source)

		if err != nil {
			return nil, err
		}

		wishlists = append(wishlists, wishlist)
	}
	return
}
