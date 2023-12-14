package main

import (
	"fmt"
	"time"
)

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	TimeStamp int64
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		TimeStamp: time.Now().UnixNano(),
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (o *Order) String() string {
	return fmt.Sprintf("[size : %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

type OrderBook struct {
	Asks []*Limit
	Bids []*Limit
}
