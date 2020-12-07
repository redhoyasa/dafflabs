package tokopedia

import (
	"context"
	"github.com/ad2games/vcr-go"
	"github.com/gocolly/colly"
	"github.com/gojektech/heimdall/v6/hystrix"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_GetItem(t *testing.T) {
	t.Run("Should get item info", func(t *testing.T) {
		vcr.Start("tkpd_get_product_200", nil)
		defer vcr.Stop()

		c := colly.NewCollector()
		c.Async = false

		client := Client{
			httpClient: hystrix.NewClient(),
		}

		expected := &product.Item{
			Name:          "Matchamu Matcha Latte 20pcs",
			OriginalPrice: 100407,
			CurrentPrice:  90366,
			Source:        "https://www.tokopedia.com/matchamu/matchamu-matcha-latte-20pcs",
		}

		item, err := client.GetItem(context.Background(), "https://www.tokopedia.com/matchamu/matchamu-matcha-latte-20pcs")
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, expected, item)
	})

	t.Run("Should return error", func(t *testing.T) {
		vcr.Start("tkpd_get_product_500", nil)
		defer vcr.Stop()

		c := colly.NewCollector()
		c.Async = false

		client := Client{
			httpClient: hystrix.NewClient(),
		}

		item, err := client.GetItem(context.Background(), "https://www.tokopedia.com/matchamu/matchamu-matcha-latte-20pcsss")
		assert.NotNil(t, err)
		assert.Nil(t, item)
		assert.Equal(t, "failed to fetch product", err.Error())
	})
}
