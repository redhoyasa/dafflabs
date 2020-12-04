package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/redhoyasa/dafflabs/internal/service/wishlist"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type wishlistSvcIFace interface {
	Add(wishlist wishlist.Wishlist) error
	FetchByCustomer(customerRefID string) ([]wishlist.Wishlist, error)
	FetchAll() ([]wishlist.Wishlist, error)
}

type Handler struct {
	svc wishlistSvcIFace
}

func NewHandler(wishlistSvc wishlistSvcIFace) *Handler {
	return &Handler{
		svc: wishlistSvc,
	}
}

func (h *Handler) AddWishlist(c echo.Context) error {
	echoRequest := c.Request()
	requestBody := echoRequest.Body
	defer requestBody.Close()

	payload, err := ioutil.ReadAll(requestBody)
	if err != nil {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(500)
		return echo.NewHTTPError(500, "error")
	}

	var wishlist wishlist.Wishlist

	if err := json.Unmarshal(payload, &wishlist); err != nil {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(500)
		return echo.NewHTTPError(500, "error")
	}

	err = h.svc.Add(wishlist)
	if err != nil {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(500)
		log.Err(err).Msg(err.Error())
		return echo.NewHTTPError(500, "error")
	}

	return c.JSON(http.StatusOK, "success")
}
