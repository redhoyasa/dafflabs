package tokopedia

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gojektech/heimdall/v6/hystrix"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"io/ioutil"
	"net/http"
)

type Client struct {
	httpClient *hystrix.Client
}

const (
	baseUrl = "https://scraper.dafflabs.workers.dev/product/tokopedia"
)

type scraperResponse struct {
	Data  tokopediaProduct `json:"data"`
	Error string           `json:"error"`
}

type tokopediaProduct struct {
	Name          string `json:"name"`
	CurrentPrice  int64  `json:"current_price"`
	OriginalPrice int64  `json:"original_price"`
	DiscountRate  int64  `json:"discount_rate"`
}

func NewClient(httpClient *hystrix.Client) (*Client, error) {
	c := new(Client)
	c.httpClient = httpClient
	return c, nil
}

func (c *Client) GetItem(ctx context.Context, source string) (item *product.Item, err error) {
	req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
	q := req.URL.Query()
	q.Add("url", source)
	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var scraperResponse scraperResponse
	_ = json.Unmarshal(body, &scraperResponse)

	if scraperResponse.Error != "" {
		return nil, errors.New("failed to fetch product")
	}
	return toItemModel(scraperResponse.Data, source), nil
}

func toItemModel(tkpdProduct tokopediaProduct, source string) (item *product.Item) {
	item = &product.Item{}
	item.Name = tkpdProduct.Name
	item.Source = source
	item.CurrentPrice = tkpdProduct.CurrentPrice
	item.OriginalPrice = tkpdProduct.OriginalPrice
	return item
}
