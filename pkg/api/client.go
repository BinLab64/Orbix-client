package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// SideType define side type of order
type SideType string

const (
	SideTypeBuy  SideType = "buy"
	SideTypeSell SideType = "sell"
)

type OrderStatusType string

const (
	OrderStatusOpen  OrderStatusType = "open"
	OrderStatusClose OrderStatusType = "close"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// Default Constants
const (
	DefaultBaseURL   = "https://satangcorp.com"
	DefaultUserAgent = "Bin64/1.0"
	DefaultTimeOut   = 10 * time.Second
)

type ClientAuth struct {
	apiKey    string
	apiSecret string
}

func NewClientAuth(apikey string, apiSecret string) ClientAuth {
	return ClientAuth{
		apiKey:    apikey,
		apiSecret: apiSecret,
	}
}

type Client struct {
	ClientAuth
	HttpClient *http.Client
	BaseURL    string
	UserAgent  string
	Logger     *slog.Logger
}

type ClientOptions struct {
	ClientAuth
	BaseURL   string
	UserAgent string
	Logger    *slog.Logger
}

func newDefaultLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func NewClient(opts ClientOptions) *Client {
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultBaseURL
	}

	if opts.UserAgent == "" {
		opts.UserAgent = DefaultUserAgent
	}

	if opts.Logger == nil {
		opts.Logger = newDefaultLogger()
	}

	return &Client{
		ClientAuth: opts.ClientAuth,
		HttpClient: http.DefaultClient,
		BaseURL:    opts.BaseURL,
		UserAgent:  opts.UserAgent,
		Logger:     opts.Logger,
	}
}

// parseRequest prepares the request and constructs the full URL.
func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// Set request options from user
	for _, opt := range opts {
		opt(r)
	}

	// Ensure the query param and form are initialized
	r.ensureInitialized()

	err = c.buildHeader(r)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error signing payload: %v", err.Error()))
		return err
	}

	c.buildFullURL(r)

	return nil
}

func (c *Client) buildHeader(r *request) error {
	header := make(http.Header)

	if r.secType == secTypeSigned && r.method == http.MethodGet {
		signature, err := Sign(c.ClientAuth.apiSecret, nil)
		if err != nil {
			return fmt.Errorf("error signing payload: %w", err)
		}
		header.Set("Authorization", "TDAX-API "+c.ClientAuth.apiKey)
		header.Set("Signature", signature)
	}

	if r.secType == secTypeSigned && (r.method == http.MethodPost || r.method == http.MethodDelete) {
		header.Set("Content-Type", "application/json")
		header.Set("Authorization", "TDAX-API "+c.ClientAuth.apiKey)
		var p params
		err := json.Unmarshal(r.bodyBuffer, &p)
		if err != nil {
			return fmt.Errorf("err building header: %w", err)
		}

		r.body = bytes.NewBuffer(r.bodyBuffer)
		signature, err := Sign(c.ClientAuth.apiSecret, p)
		if err != nil {
			return fmt.Errorf("error signing payload: %w", err)
		}
		header.Set("Signature", signature)
	}

	r.header = header
	return nil
}

// Constructs the full URL for the request.
func (c *Client) buildFullURL(r *request) {
	var sb strings.Builder
	sb.WriteString(c.BaseURL)
	sb.WriteString(r.endpoint)

	// Encode the query parameters
	queryString := r.query.Encode()
	if queryString != "" {
		sb.WriteString("?")
		sb.WriteString(queryString)
	}
	r.fullURL = sb.String()
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, r.fullURL, r.body)
	req.Header = r.header

	fmt.Printf("\n\n<request>\n%+v\n</request>\n\n", req)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.Logger.Debug(
		"Orbix API Call",
		slog.String("method", r.method),
		slog.String("header", fmt.Sprintf("%+v", r.header)),
		slog.String("url", r.fullURL),
	)

	c.Logger.Debug(
		"Orbix API Reponse",
		slog.String("status", res.Status),
		slog.String("header", fmt.Sprintf("%+v", res.Header)),
		slog.String("body", string(data)),
	)

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("API request failed with status %d: %s", res.StatusCode, string(data))
	}
	return data, nil
}

func (c *Client) SetBaseURL(url string) *Client {
	c.BaseURL = url
	return c
}

//
//
//
//
//
//

// Get balances and addresses

// Get exchange information
// *
// /api/v3/exchangeInfo
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

// *
// /api/v3/depth
func (c *Client) NewOrderbookDepthService(symbol string) *OrderbookDepthService {
	return &OrderbookDepthService{
		c:      c,
		symbol: symbol,
	}
}

// /api/v3/klines
func (c *Client) NewKlineService() *KlineService {
	return &KlineService{c: c}
}

//
// ---------------------------------------------------------------------------
//  Order Services

// *
// /api/orders/
func (c *Client) NewOrderbookService(pair string) *OrderbookService {
	return &OrderbookService{c: c, pair: pair}
}

// *
// /api/orders/<Order Id>
func (c *Client) NewGetOrderByIdService(orderId string, pair string) *GetOrderByIdService {
	return &GetOrderByIdService{c: c, pair: pair, orderId: orderId}
}

// /api/orderbook-tickers/
func (c *Client) NewOrderbookTickerService() *OrderbookTickerService {
	return &OrderbookTickerService{c: c}
}

// *
// /api/orders/

func (c *Client) NewCreateOrderService(pair string, side SideType, orderType OrderType, price string, amount string) *CreateOrderService {
	return &CreateOrderService{c: c, pair: pair, side: side, orderType: orderType, price: price, amount: amount}
}

// *
// /api/orders/user

func (c *Client) NewListCurrentOrdersService(pair string, limit int, offset int) *ListCurrentOrdersService {
	return &ListCurrentOrdersService{c: c, pair: pair, limit: limit, offset: offset}
}

// /api/orders/<Order Id>
func (c *Client) NewCancelOrderService(orderId string, pair string) *CancelOrderService {
	return &CancelOrderService{c: c, orderId: orderId, pair: pair}
}

// /api/orders/all
func (c *Client) NewCancelAllOrdersService(pair string) *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c, pair: pair}
}

// Get 24 hrs. ticker
// /api/v3/ticker/24hr
func (c *Client) NewList24HrPriceChangeStatsService() *List24HrPriceChangeStatsService {
	return &List24HrPriceChangeStatsService{c: c}
}

// *
// GET Get balances and addresses
// /api/users/me
func (c *Client) NewListBalanceAddressService() *ListBalanceAddressService {
	return &ListBalanceAddressService{c: c}
}

// Get aggregate trade
// /api/v3/aggTrades
func (c *Client) NewAggregateTradeService() *AggregateTradeService {
	return &AggregateTradeService{c: c}
}

// GET Ping -- Get all configs
// /api/v3/ping
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// POST Generate a listen key (UserStream) Create a new listen key. The listen key will be expired after 60 minutes
// /api/v3/userDataStream
func (c *Client) NewListenKeyService() *ListenKeyService {
	return &ListenKeyService{c: c}
}

// POST Keep-alive a listen key (UserStream) Keep-alive a listen key for 30 minutes
// /api/v3/userDataStream
func (c *Client) NewKeepAliveListenKeySerice() *KeepAliveListenKeySerice {
	return &KeepAliveListenKeySerice{c: c}
}

// GET Fiat deposit histories
// /api/bank-account-deposits
func (c *Client) NewFiatDepositHistoryService() *FiatDepositHistoryService {
	return &FiatDepositHistoryService{c: c}
}

// GET fiat histories
// /api/fiat-withdrawals
func (c *Client) NewFiatWithdrawalHistoryService() *FiatWithdrawalHistoryService {
	return &FiatWithdrawalHistoryService{c: c}
}

// GET crypto deposit history
// /api/crypto-deposits
func (c *Client) NewCryptoDepositHistoryService() *CryptoDepositHistoryService {
	return &CryptoDepositHistoryService{c: c}
}

// GET histories Required permission: withdrawal_list
// /api/crypto-withdrawals
func (c *Client) NewCryptoWithdrawalHistoryService() *CryptoWithdrawalHistoryService {
	return &CryptoWithdrawalHistoryService{c: c}
}

// GET trade history
// /api/trade-history
func (c *Client) NewTradeHistoryService() *TradeHistoryService {
	return &TradeHistoryService{c: c}
}

// GET configs Get all configs
// /api/configs
func (c *Client) NewAllConfigsService() *AllConfigsService {
	return &AllConfigsService{c: c}
}
