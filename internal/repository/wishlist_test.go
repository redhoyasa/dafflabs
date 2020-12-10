package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redhoyasa/dafflabs/internal/database"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbMock database.Database
var sqlMock sqlmock.Sqlmock

func TestWishlistRepo_Insert(t *testing.T) {

	t.Run("Should insert row", func(t *testing.T) {
		dbMock, sqlMock, _ = sqlmock.New()
		defer dbMock.Close()

		sqlMock.ExpectBegin()
		sqlMock.
			ExpectExec("INSERT INTO wishes").
			WithArgs("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "12342", "PS 5", 100, 1000, "amazon.com").
			WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()

		r := &WishRepo{
			db: dbMock,
		}

		w := wishlist.Wish{
			WishID:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			CustomerRefID: "12342",
			ProductName:   "PS 5",
			CurrentPrice:  100,
			OriginalPrice: 1000,
			Source:        "amazon.com",
		}

		err := r.Insert(w)

		assert.Nil(t, err)
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}

func TestWishlistRepo_FetchByCustomer(t *testing.T) {
	t.Run("Should fetch rows", func(t *testing.T) {
		dbMock, sqlMock, _ = sqlmock.New()
		defer dbMock.Close()

		expectedQuery := `
			SELECT
				id,
				customer_ref_id,
				product_name,
				current_price,
				original_price,
				source
			FROM wishes
			WHERE
				is_deleted = 'false'
				AND customer_ref_id = ?
		`

		sqlMock.
			ExpectQuery(expectedQuery).
			WithArgs("customer1").
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "customer_ref_id", "product_name", "current_price", "original_price", "source"}).
				AddRow(1, "customer1", "PS 5", 100, 1000, "amazon.com").
				AddRow(2, "customer1", "Ferrari", 1000, 10000, "ferrari.com"))

		r := &WishRepo{
			db: dbMock,
		}

		wishes, err := r.FetchByCustomer("customer1")

		assert.Nil(t, err)
		assert.Equal(t, 2, len(wishes))
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestWishlistRepo_FetchAll(t *testing.T) {
	t.Run("Should fetch all rows", func(t *testing.T) {
		dbMock, sqlMock, _ = sqlmock.New()
		defer dbMock.Close()

		expectedQuery := `
			SELECT
				id,
				customer_ref_id,
				product_name,
				current_price,
				original_price,
				source
			FROM wishes
			WHERE
				is_deleted = 'false'
		`

		sqlMock.
			ExpectQuery(expectedQuery).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "customer_ref_id", "product_name", "current_price", "original_price", "source"}).
				AddRow(1, "customer1", "PS 5", 100, 1000, "amazon.com").
				AddRow(2, "customer2", "Ferrari", 1000, 10000, "ferrari.com"))

		r := &WishRepo{
			db: dbMock,
		}

		wishes, err := r.FetchAll()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(wishes))
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
