package repository

import (
	"context"
	"fmt"
	"github.com/redhoyasa/dafflabs/internal/database"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
)

type WishRepo struct {
	db database.Database
}

func NewWishRepo(db database.Database) *WishRepo {
	return &WishRepo{
		db: db,
	}
}

func (w *WishRepo) Insert(wish wishlist.Wish) (err error) {
	query := fmt.Sprintf(
		`INSERT INTO wishes(
					id,
					customer_ref_id, 
					product_name, 
					current_price, 
					original_price,
					discount_rate,
					source
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7
				)`)

	tx, err := w.db.Begin()
	if err != nil {
		return
	}

	if _, err = tx.Exec(query, wish.WishID, wish.CustomerRefID, wish.ProductName, wish.CurrentPrice, wish.OriginalPrice, wish.DiscountRate, wish.Source); err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()
	return
}

func (w *WishRepo) FetchByCustomer(customerRefID string) (wishes []wishlist.Wish, err error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			customer_ref_id,
			product_name,
			current_price,
			original_price,
			discount_rate,
			source,
			updated_at
		FROM wishes
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
		wish := wishlist.Wish{}
		err := rows.Scan(
			&wish.WishID,
			&wish.CustomerRefID,
			&wish.ProductName,
			&wish.CurrentPrice,
			&wish.OriginalPrice,
			&wish.DiscountRate,
			&wish.Source,
			&wish.UpdatedAt)

		if err != nil {
			return nil, err
		}

		wishes = append(wishes, wish)
	}
	return
}

func (w *WishRepo) FetchAll() (wishes []wishlist.Wish, err error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			customer_ref_id,
			product_name,
			current_price,
			original_price,
			discount_rate,
			source,
			updated_at
		FROM wishes
		WHERE
			is_deleted = 'false'
	`)

	rows, err := w.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		wish := wishlist.Wish{}
		err := rows.Scan(
			&wish.WishID,
			&wish.CustomerRefID,
			&wish.ProductName,
			&wish.CurrentPrice,
			&wish.OriginalPrice,
			&wish.DiscountRate,
			&wish.Source,
			&wish.UpdatedAt)

		if err != nil {
			return nil, err
		}

		wishes = append(wishes, wish)
	}
	return
}

func (w *WishRepo) Fetch(wishID string) (wish *wishlist.Wish, err error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			customer_ref_id,
			product_name,
			current_price,
			original_price,
			discount_rate,
			source
		FROM wishes
		WHERE
			id = $1
	`)

	row := w.db.QueryRowContext(context.Background(), query, wishID)
	if err != nil {
		return nil, err
	}

	wish = &wishlist.Wish{}
	if err = row.Scan(&wish.WishID, &wish.CustomerRefID, &wish.ProductName, &wish.CurrentPrice, &wish.OriginalPrice, &wish.DiscountRate, &wish.Source); err != nil {
		return nil, err
	}

	return
}

func (w *WishRepo) Delete(wishID string) (err error) {
	query := fmt.Sprintf("DELETE FROM wishes WHERE id = $1")

	tx, err := w.db.Begin()
	if err != nil {
		return
	}

	if _, err = tx.Exec(query, wishID); err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()

	return nil
}

func (w *WishRepo) Update(wishID string, originalPrice, currentPrice int64, discountRate float64) (err error) {
	query := fmt.Sprintf(`
			UPDATE wishes 
			SET
				original_price = $2,
				current_price = $3,
				discount_rate = $4,
				updated_at = NOW()
			WHERE 
				id = $1
			`)

	tx, err := w.db.Begin()
	if err != nil {
		return
	}

	if _, err = tx.Exec(query, wishID, originalPrice, currentPrice, discountRate); err != nil {
		_ = tx.Rollback()
		return
	}

	_ = tx.Commit()

	return nil
}
