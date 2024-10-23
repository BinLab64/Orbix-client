package api

import (
	"context"
	"encoding/json"
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
