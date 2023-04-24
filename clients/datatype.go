package clients

import (
	"fmt"
	"time"
)

type TokenPair struct {
	BaseToken  string `json:"base_token"`
	QuoteToken string `json:"quote_token"`
}

type Quantity struct {
	Token string  `json:"token"`
	Value float64 `json:"value,string"`
}

type QuoteRequest struct {
	TokenPair     TokenPair `json:"token_pair"`
	Quantity      Quantity  `json:"quantity"`
	Side          string    `json:"side"`
	ClientOrderId string    `json:"client_order_id"`
}

type OrderRequest struct {
	TokenPair     TokenPair `json:"token_pair"`
	Quantity      Quantity  `json:"quantity"`
	Side          string    `json:"side"`
	OrderType     string    `json:"order_type"`
	TimeInForce   string    `json:"time_in_force"`
	LimitPrice    float64   `json:"limit_price"`
	SlippageBps   float64   `json:"slippage_bps"`
	ClientOrderId string    `json:"client_order_id"`
}

type QuoteExecutionRequest struct {
	FxQuoteId string `json:"fx_quote_id"`
	Side      string `json:"side"`
}

type FalconXError struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

type FalconXWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Side    string `json:"side"`
}

type QuoteResponse struct {
	Status        string           `json:"status"`
	FxQuoteId     string           `json:"fx_quote_id"`
	BuyPrice      float64          `json:"buy_price,string"`
	SellPrice     float64          `json:"sell_price,string"`
	Platform      string           `json:"platform"`
	TokenPair     TokenPair        `json:"token_pair"`
	Quantity      Quantity         `json:"quantity_requested"`
	PositionIn    Quantity         `json:"position_in"`
	PositionOut   Quantity         `json:"position_out"`
	SideRequested string           `json:"side_requested"`
	QuoteTime     time.Time        `json:"t_quote"`
	ExpiryTime    time.Time        `json:"t_expiry"`
	ExecutionTime time.Time        `json:"t_execute"`
	IsFilled      bool             `json:"is_filled"`
	TraderEmail   string           `json:"trader_email"`
	Error         FalconXError     `json:"error"`
	Warnings      []FalconXWarning `json:"warnings"`
	ClientOrderId string           `json:"client_order_id"`
}

type OrderResponse struct {
	Status        string           `json:"status"`
	FxQuoteId     string           `json:"fx_quote_id"`
	BuyPrice      float64          `json:"buy_price,string"`
	SellPrice     float64          `json:"sell_price,string"`
	Platform      string           `json:"platform"`
	TokenPair     TokenPair        `json:"token_pair"`
	Quantity      Quantity         `json:"quantity_requested"`
	SideRequested string           `json:"side_requested"`
	QuoteTime     time.Time        `json:"t_quote"`
	ExpiryTime    time.Time        `json:"t_expiry"`
	ExecutionTime time.Time        `json:"t_execute"`
	IsFilled      bool             `json:"is_filled"`
	GrossFeeBps   float64          `json:"gross_fee_bps,string"`
	GrossFeeUSD   float64          `json:"gross_fee_usd,string"`
	RebateBps     float64          `json:"rebate_bps,string"`
	RebateUSD     float64          `json:"rebate_usd,string"`
	FeeBps        float64          `json:"fee_bps,string"`
	FeeUSD        float64          `json:"fee_usd,string"`
	SideExecuted  string           `json:"side_executed"`
	TraderEmail   string           `json:"trader_email"`
	OrderType     string           `json:"order_type"`
	TimeInForce   string           `json:"time_in_force"`
	LimitPrice    float64          `json:"limit_price,string"`
	SlippageBps   float64          `json:"slippage_bps,string"`
	Error         FalconXError     `json:"error"`
	Warnings      []FalconXWarning `json:"warnings"`
	ClientOrderId string           `json:"client_order_id"`
}

type Balance struct {
	Token    string  `json:"token"`
	Balance  float64 `json:"balance"`
	Platform string  `json:"platform"`
}

type TotalBalance struct {
	Token        string  `json:"token"`
	TotalBalance float64 `json:"total_balance"`
}

type Transfer struct {
	Type       string    `json:"type"`
	Platform   string    `json:"platform"`
	Token      string    `json:"token"`
	Quantity   float64   `json:"quantity,string"`
	CreateTime time.Time `json:"t_create"`
	Status     string    `json:"status"`
}

type TradeVolume struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	USDVolume float64   `json:"usd_volume"`
}

type TradeLimit struct {
	Available float64 `json:"available"`
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
}
type TradeLimits struct {
	GrossLimits TradeLimit `json:"gross_limits"`
	NetLimits   TradeLimit `json:"net_limits"`
}

type TradeSizeLimit struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}
type TradeSize struct {
	Platform                 string         `json:"platform"`
	TokenPair                TokenPair      `json:"token_pair"`
	TradeSizeLimitQuoteToken TradeSizeLimit `json:"trade_size_limits_in_quote_token"`
}

type SubscriptionRequest struct {
	TokenPair       TokenPair `json:"token_pair"`
	Quantity        []float64 `json:"quantity"`
	ClientRequestID string    `json:"client_request_id"`
}

type UserConfigRequest struct {
	MessageType     string `json:"message_type"`
	ClientRequestID string `json:"client_request_id"`
}

type UserConfigResponse struct {
	MessageType     string      `json:"message_type"`
	ClientRequestID string      `json:"client_request_id"`
	Success         bool        `json:"success"`
	Data            interface{} `json:"data"`
	Error           interface{} `json:"error"`
}

type Error struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Status Code: %d \n Reason: %s", e.Code, e.Reason)
}
