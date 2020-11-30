package tokopedia

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/redhoyasa/dafflabs/internal/repository/product"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	scraper *colly.Collector
}

type tokopediaProduct struct {
	Name  string `selector:"h1[data-testid=lblPDPDetailProductName]"`
	Price string `selector:"h3[data-testid=lblPDPDetailProductPrice]"`
}

func NewClient(scraper *colly.Collector) (*Client, error){
	c := new(Client)
	c.scraper = scraper
	return c, nil
}

func (c *Client) GetItem(ctx context.Context, source string) (item product.Item, err error) {
	c.scraper.OnHTML("div[data-testid=pdpContainer]", func(e *colly.HTMLElement) {
		product := &tokopediaProduct{}
		e.Unmarshal(product)

		err = toItemModel(product, &item, source)
		if err != nil {
			return
		}
	})

	c.scraper.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err = c.scraper.Request("GET", source, nil, nil, getHeader())
	if err != nil {
		return
	}

	return
}

func toItemModel(tkpdProduct *tokopediaProduct, item *product.Item, source string) error {
	tkpdProduct.Price = strings.Replace(tkpdProduct.Price, "Rp", "", -1)
	tkpdProduct.Price = strings.Replace(tkpdProduct.Price, ".", "", -1)
	parsedPrice, err := strconv.ParseInt(tkpdProduct.Price, 10, 64)
	if err != nil {
		return err
	}

	item.Name = tkpdProduct.Name
	item.Source = source
	item.Price = parsedPrice
	return nil
}

func getHeader() http.Header {
	header := http.Header{}
	header.Add("Accept", "*/*")
	header.Add("Upgrade-Insecure-Requests", "1")
	header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	return header
}