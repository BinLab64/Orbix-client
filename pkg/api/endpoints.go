package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ExchangeInfoService struct {
	c *Client
}

type ExchangeInfo struct {
	Timezone        string               `json:"timezone"`
	ServerTime      int64                `json:"serverTime"`
	RateLimits      any                  `json:"rateLimits"`      // If rateLimits are not defined, use interface{}
	ExchangeFilters any                  `json:"exchangeFilters"` // If exchangeFilters are empty
	Symbols         []ExchangeInfoSymbol `json:"symbols"`
}

type ExchangeInfoSymbol struct {
	Symbol                     string               `json:"symbol"`
	Status                     string               `json:"status"`
	BaseAsset                  string               `json:"baseAsset"`
	BaseAssetPrecision         int                  `json:"baseAssetPrecision"`
	QuoteAsset                 string               `json:"quoteAsset"`
	QuotePrecision             int                  `json:"quotePrecision"`
	BaseCommissionPrecision    int                  `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int                  `json:"quoteCommissionPrecision"`
	OrderTypes                 []string             `json:"orderTypes"`
	IcebergAllowed             bool                 `json:"icebergAllowed"`
	OcoAllowed                 bool                 `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                 `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool                 `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                 `json:"isMarginTradingAllowed"`
	Filters                    []ExchangeInfoFilter `json:"filters"`
}

type ExchangeInfoFilter struct {
	FilterType string `json:"filterType"`
	TickSize   string `json:"tickSize"` // Assuming TickSize as a string
}

func (s *ExchangeInfoService) Do(ctx context.Context, opt ...RequestOption) (exchangeInfo *ExchangeInfo, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/exchangeInfo",
		secType:  secTypeNone,
	}

	data, err := s.c.callAPI(ctx, r, opt...)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &exchangeInfo); err != nil {
		return nil, err
	}
	return exchangeInfo, nil
}

type OrderbookDepthService struct {
	c      *Client
	symbol string
	limit  *int
}

type OrderbookDepth struct {
	LastUpdateId int64       `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

func (s *OrderbookDepthService) Limit(limit int) *OrderbookDepthService {
	s.limit = &limit
	return s
}

func (s *OrderbookDepthService) Do(ctx context.Context, opt ...RequestOption) (orderbookDepth *OrderbookDepth, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/v3/depth",
		secType:  secTypeNone,
	}
	r.setQueryParam("symbol", s.symbol)

	if s.limit != nil {
		if *s.limit < 5 || *s.limit > 5000 {
			return nil, fmt.Errorf("error: invalid limit [%v], must be between 5 and 5000", *s.limit)
		}
		r.setQueryParam("limit", *s.limit)
	}

	data, err := s.c.callAPI(ctx, r, opt...)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &orderbookDepth); err != nil {
		return nil, err
	}

	return orderbookDepth, nil
}

type KlineService struct {
	c *Client
}

// Get 24 hrs. ticker
// /api/v3/ticker/24hr
type List24HrPriceChangeStatsService struct {
	c *Client
}

// GET Get balances and addresses
type ListBalanceAddressService struct {
	c *Client
}

type CoinSymbol string

// Main struct for the user data
type User struct {
	ID                        int                   `json:"id"`
	Email                     string                `json:"email"`
	IdentityVerificationLevel string                `json:"identity_verification_level"`
	TFAEnabled                []string              `json:"tfa_enabled"`
	APIKeys                   []APIKey              `json:"api_keys"`
	AntiPhishingCode          string                `json:"anti_phishing_code"`
	IsAuthorizedDevice        bool                  `json:"is_authorized_device"`
	Wallets                   map[CoinSymbol]Wallet `json:"wallets"`
}

// Struct for API key details
type APIKey struct {
	APIKey      string   `json:"APIKey"`
	Label       string   `json:"Label"`
	Status      int      `json:"Status"`
	Permissions []string `json:"Permissions"`
	CreatedAt   string   `json:"CreatedAt"`
}

// Struct for wallet details
type Wallet struct {
	Addresses        []Address `json:"addresses"`
	AvailableBalance string    `json:"available_balance"`
}

// Struct for address details
type Address struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
	Network string `json:"network"`
}

func (s *ListBalanceAddressService) Do(ctx context.Context, opt ...RequestOption) (user *User, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/api/users/me",
		secType:  secTypeSigned,
	}

	data, err := s.c.callAPI(ctx, r, opt...)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// Get aggregate trade
// /api/v3/aggTrades
type AggregateTradeService struct {
	c *Client
}

// GET Ping -- Get all configs
// /api/v3/ping
type PingService struct {
	c *Client
}

// POST Generate a listen key (UserStream) Create a new listen key. The listen key will be expired after 60 minutes
// /api/v3/userDataStream
type ListenKeyService struct {
	c *Client
}

// POST Keep-alive a listen key (UserStream) Keep-alive a listen key for 30 minutes
// /api/v3/userDataStream
type KeepAliveListenKeySerice struct {
	c *Client
}

// GET Fiat deposit histories
// /api/bank-account-deposits
type FiatDepositHistoryService struct {
	c *Client
}

// GET fiat histories
// /api/fiat-withdrawals
type FiatWithdrawalHistoryService struct {
	c *Client
}

// GET crypto deposit history
// /api/crypto-deposits
type CryptoDepositHistoryService struct {
	c *Client
}

// GET histories Required permission: withdrawal_list
// /api/crypto-withdrawals
type CryptoWithdrawalHistoryService struct {
	c *Client
}

// GET trade history
// /api/trade-history
type TradeHistoryService struct {
	c *Client
}

// GET configs Get all configs
// /api/configs
type AllConfigsService struct {
	c *Client
}
