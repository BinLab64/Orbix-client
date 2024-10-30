package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GET Get order book

type OrderbookService struct {
	c    *Client
	pair string
	side *SideType
}

type OrderbookItem struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type PartialOrderbook []OrderbookItem

type Orderbook struct {
	Bids []OrderbookItem `json:"bid"`
	Asks []OrderbookItem `json:"ask"`
}

func (s *OrderbookService) Side(side SideType) *OrderbookService {
	s.side = &side
	return s
}

func (s *OrderbookService) Do(ctx context.Context, opts ...RequestOption) (ob *Orderbook, err error) {
	var hasSideParam bool
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/orders/",
		secType:  secTypeNone,
	}
	r.setQueryParam("pair", s.pair)
	if s.side != nil {
		r.setQueryParam("side", *s.side)
		hasSideParam = true
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if hasSideParam {
		partialOrderbook := make([]OrderbookItem, 0, 10)
		if err := json.Unmarshal(data, &partialOrderbook); err != nil {
			return nil, err
		}
		if *s.side == SideTypeBuy {
			return &Orderbook{Bids: partialOrderbook}, nil
		}
		return &Orderbook{Asks: partialOrderbook}, nil
	}

	if err := json.Unmarshal(data, &ob); err != nil {
		return nil, err
	}
	return ob, nil
}

// GET Get order by ID

type GetOrderByIdService struct {
	c       *Client
	pair    string
	orderId string
}

type Order struct {
	ID              int       `json:"id"`
	Type            OrderType `json:"type"`
	Price           string    `json:"price"`
	Amount          string    `json:"amount"`
	RemainingAmount string    `json:"remaining_amount"`
	AveragePrice    string    `json:"average_price"`
	Side            SideType  `json:"side"`
	Cost            string    `json:"cost"`
	CreatedAt       string    `json:"created_at"`
	Status          string    `json:"status"`
}

func (s *GetOrderByIdService) Do(ctx context.Context, opts ...RequestOption) (order *Order, err error) {
	// /api/orders/1234000?pair=usdt_thb
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/orders/" + s.orderId,
		secType:  secTypeNone,
	}
	r.setQueryParam("pair", s.pair)

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return order, nil
}

// GET Get order book ticker
type OrderbookTickerService struct {
	c *Client
}

// POST Create order
type CreateOrderService struct {
	c         *Client
	pair      string
	side      SideType
	orderType OrderType
	price     string
	amount    string
}

type CreateOrderRequestBody struct {
	Amount string    `json:"amount"`
	Nonce  int64     `json:"nonce"`
	Pair   string    `json:"pair"`
	Price  string    `json:"price"`
	Side   SideType  `json:"side"`
	Type   OrderType `json:"type"`
}

func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (order *Order, err error) {

	body, err := json.Marshal(CreateOrderRequestBody{
		Amount: s.amount,
		// Nonce:  fmt.Sprintf("%v", time.Now().UnixMilli()),
		// Nonce: time.Now().UnixMilli(),
		Nonce: 2731832,
		Pair:  s.pair,
		Price: s.price,
		Side:  s.side,
		Type:  s.orderType,
	})
	if err != nil {
		return nil, fmt.Errorf("err marshalling JSON: %w", err)
	}

	r := &request{
		method:     http.MethodPost,
		endpoint:   "/api/orders/",
		secType:    secTypeSigned,
		bodyBuffer: body,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return order, nil
}

// GET List your orders
type ListCurrentOrdersService struct {
	c      *Client
	pair   string
	limit  int
	offset int
	status *OrderStatusType
	side   *SideType
}

func (s *ListCurrentOrdersService) Status(status OrderStatusType) *ListCurrentOrdersService {
	s.status = &status
	return s
}

func (s *ListCurrentOrdersService) Side(side SideType) *ListCurrentOrdersService {
	s.side = &side
	return s
}

func (s *ListCurrentOrdersService) Do(ctx context.Context, opts ...RequestOption) (orders *[]Order, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/orders/user",
		secType:  secTypeSigned,
	}
	r.setQueryParams(params{
		"pair":   s.pair,
		"limit":  s.limit,
		"offset": s.offset,
	})
	if s.side != nil {
		r.addQueryParam("side", *s.side)
	}
	if s.status != nil {
		r.addQueryParam("status", *s.status)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// DEL Cancel order
type CancelOrderService struct {
	c       *Client
	orderId string
	pair    string
}

type CancelOrderRequestBody struct {
	Pair string `json:"pair"`
}

func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (err error) {

	body, err := json.Marshal(CancelOrderRequestBody{
		Pair: s.pair,
	})
	r := &request{
		method:     http.MethodDelete,
		endpoint:   "/api/orders/" + s.orderId,
		secType:    secTypeSigned,
		bodyBuffer: body,
	}
	if err != nil {
		return fmt.Errorf("err marshalling JSON: %w", err)
	}
	_, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// DEL Cancel all order
type CancelAllOrdersService struct {
	c    *Client
	pair string
}

type CancelAllOrdersResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (s *CancelAllOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *[]CancelAllOrdersResponse, err error) {

	body, err := json.Marshal(CancelOrderRequestBody{
		Pair: s.pair,
	})
	r := &request{
		method:     http.MethodDelete,
		endpoint:   "/api/orders/all",
		secType:    secTypeSigned,
		bodyBuffer: body,
	}
	if err != nil {
		return nil, fmt.Errorf("err marshalling JSON: %w", err)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
