package main

import (
	"encoding/json"
	"github/tetrex/crypto_exchange/orderbook"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)

	e.Start(":3000")
}

// ------------------
type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

// ------------------------

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)

type PlaceOrderRequest struct {
	Type   OrderType //limit or market
	Bid    bool
	Size   float64
	Price  float64
	Market Market
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	ob := ex.orderbooks[market]
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	if placeOrderData.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(http.StatusOK, map[string]any{"msg": "limit order placed"})
	}

	if placeOrderData.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		return c.JSON(http.StatusOK, map[string]any{"matches": len(matches)})
	}

	return nil
}

type Order struct {
	Price     float64
	Size      float64
	Bid       bool
	TimeStamp int64
}
type OrderBookData struct {
	Asks []*Order
	Bids []*Order
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))

	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderBookData{
		Asks: []*Order{},
		Bids: []*Order{},
	}

	for _, limits := range ob.Asks() {
		for _, orders := range limits.Orders {
			o := Order{
				Price:     orders.Limit.Price,
				Size:      orders.Size,
				Bid:       orders.Bid,
				TimeStamp: orders.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limits := range ob.Bids() {
		for _, orders := range limits.Orders {
			o := Order{
				Price:     orders.Limit.Price,
				Size:      orders.Size,
				Bid:       orders.Bid,
				TimeStamp: orders.Timestamp,
			}
			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}
