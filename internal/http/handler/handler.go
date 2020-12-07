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
	Add(wishlist *wishlist.Wishlist) error
	FetchByCustomer(customerRefID string) ([]wishlist.Wishlist, error)
	FetchAll() ([]wishlist.Wishlist, error)
}

type Handler struct {
	svc wishlistSvcIFace
}

type response struct {
	httpCode int
	Data     interface{} `json:"data"`
	Error    *string     `json:"error"`
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

	var wishlist wishlist.Wishlist
	resp := response{}

	payload, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.Err(err).Msg(err.Error())
		errString := err.Error()
		resp.Error = &errString
		resp.httpCode = http.StatusInternalServerError
		return h.writeResponse(c, resp)
	}

	if err := json.Unmarshal(payload, &wishlist); err != nil {
		log.Err(err).Msg(err.Error())
		errString := err.Error()
		resp.Error = &errString
		resp.httpCode = http.StatusInternalServerError
		return h.writeResponse(c, resp)
	}

	err = h.svc.Add(&wishlist)
	if err != nil {
		log.Err(err).Msg(err.Error())
		errString := err.Error()
		resp.Error = &errString
		resp.httpCode = http.StatusInternalServerError
		return h.writeResponse(c, resp)
	}

	resp.Data = wishlist
	resp.httpCode = http.StatusOK

	return h.writeResponse(c, resp)
}

func (h *Handler) FetchCustomerWishlist(c echo.Context) error {
	echoRequest := c.Request()
	requestBody := echoRequest.Body
	defer requestBody.Close()
	customerRefID := c.QueryParam("customer_ref_id")

	resp := response{}

	wishlist, err := h.svc.FetchByCustomer(customerRefID)
	if err != nil {
		log.Err(err).Msg(err.Error())
		errString := err.Error()
		resp.Error = &errString
		resp.httpCode = http.StatusInternalServerError
		return h.writeResponse(c, resp)
	}

	resp.Data = wishlist
	resp.httpCode = http.StatusOK

	return c.JSON(resp.httpCode, resp)
}

func (h *Handler) writeResponse(c echo.Context, resp response) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	_, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return echo.NewHTTPError(resp.httpCode, resp)
	}

	return c.JSON(resp.httpCode, resp)
}
