package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type RestClient struct {
	Config     RestClientConfig
	HTTPClient *http.Client
}

type RestClientConfig struct {
	BaseURL    string
	Secret     string
	APIKey     string
	Passphrase string
}

func NewRestClient(config RestClientConfig) *RestClient {
	client := RestClient{
		Config: config,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	return &client
}

func (client *RestClient) Request(method string, url string,
	params interface{}, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", client.Config.BaseURL, url)
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	dataString := ""
	if len(data) > 0 {
		dataString = string(data)
	}

	h, err := client.Headers(method, url, timestamp, dataString)
	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = client.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		defer res.Body.Close()
		errorResponse := Error{Code: res.StatusCode}

		switch errorResponse.Code {
		case 400:
			errorResponse.Reason = "Bad Request – Invalid request format"
		case 401:
			errorResponse.Reason = "Unauthorized – Invalid API Key"
		case 403:
			errorResponse.Reason = "Forbidden – You do not have access to the requested resource"
		case 404:
			errorResponse.Reason = "Resource Not Found"
		case 500:
			errorResponse.Reason = "Internal Server Error – We had a problem with our server"
		case 503:
			errorResponse.Reason = "Service Unavailable"
		case 504:
			errorResponse.Reason = "Gateway Timeout"
		default:
			errorResponse.Reason = "Unknown Error Occured"
		}

		return res, error(errorResponse)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(&result); err != nil {
			return res, err
		}
	}

	return res, err
}

// Headers generates a map that can be used as headers to authenticate a request
func (client *RestClient) Headers(method, url, timestamp, data string) (map[string]string, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["FX-ACCESS-KEY"] = client.Config.APIKey
	header["FX-ACCESS-PASSPHRASE"] = client.Config.Passphrase
	header["FX-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := GenerateSig(message, client.Config.Secret)
	if err != nil {
		return nil, err
	}
	header["FX-ACCESS-SIGN"] = sig
	return header, nil
}

// GetTradingPairs gets a list of trading pairs you are eligible to trade
// Example: [{'base_token': 'BTC', 'quote_token': 'USD'}, {'base_token': 'ETH', 'quote_token': 'USD'}]
func (client *RestClient) GetTradingPairs() ([]TokenPair, error) {
	var result []TokenPair
	_, err := client.Request("GET", "/v1/pairs", nil, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetQuote gets a two_way, buy or sell quote for a token pair.
//         :param base: (str) base token e.g. BTC, ETH
//         :param quote: (str) quote token e.g. USD, BTC
//         :param quantity: (float, Decimal)
//         :param side: (str) 'two_way', 'buy', 'sell'
//		   Example:
//             {
//               "status": "success",
//               "fx_quote_id": "00c884b056f949338788dfb59e495377",
//               "buy_price": 12650,
//               "sell_price": null,
//               "token_pair": {
//                 "base_token": "BTC",
//                 "quote_token": "USD"
//               },
//               "quantity_requested": {
//                 "token": "BTC",
//                 "value": "10"
//               },
//               "side_requested": "buy",
//               "t_quote": "2019-06-27T11:59:21.875725+00:00",
//               "t_expiry": "2019-06-27T11:59:22.875725+00:00",
//               "is_filled": false,
//               "side_executed": null,
//               "price_executed": null,
//               "t_execute": null,
//               "client_order_id": "d6f3e1fa-e148-4009-9c07-a87f9ae78d1a"
//             }
func (client *RestClient) GetQuote(quoteParams QuoteRequest) (QuoteResponse, error) {
	var result QuoteResponse
	_, err := client.Request("POST", "/v1/quotes", quoteParams, &result)

	return result, err
}

// PlaceOrder gets a two_way, buy or sell quote for a token pair.
//         :param base: (str) base token e.g. BTC, ETH
//         :param quote: (str) quote token e.g. USD, BTC
//         :param quantity: (float, Decimal)
//         :param side: (str) 'buy', 'sell'
//         :param order_type: (str) 'market', 'limit'
//         :param time_in_force: (str) 'fok' [only required for limit orders]
//         :param limit_price: (float, Decimal) [only required for limit orders]
//         :param slippage_bps: (float, Decimal) [only valid for fok limit orders]
//		   Example:
//             {
//                 "status": "success",
//                 "fx_quote_id": "00c884b056f949338788dfb59e495377",
//                 "buy_price": 8545.12,
//                 "sell_price": null,
//                 "platform": "api",
//                 "token_pair": {
//                     "base_token": "BTC",
//                     "quote_token": "USD"
//                 },
//                 "quantity_requested": {
//                     "token": "BTC",
//                     "value": "10"
//                 },
//                 "side_requested": "buy",
//                 "t_quote": "2019-06-27T11:59:21.875725+00:00",
//                 "t_expiry": "2019-06-27T11:59:22.875725+00:00",
//                 "is_filled": true,
//                 "gross_fee_bps": 8,
//                 "gross_fee_usd": 101.20,
//                 "rebate_bps": 3,
//                 "rebate_usd": 37.95,
//                 "fee_bps": 5,
//                 "fee_usd": 63.25,
//                 "side_executed": "buy",
//                 "trader_email": "trader@company.com",
//                 "order_type": "limit",
//                 "time_in_force": "fok",
//                 "limit_price": 8547.11,
//                 "slippage_bps": 2,
//                 "error": null,
//                 "client_order_id": "d6f3e1fa-e148-4009-9c07-a87f9ae78d1a"
//             }
func (client *RestClient) PlaceOrder(orderParams OrderRequest) (OrderResponse, error) {
	var result OrderResponse
	_, err := client.Request("POST", "/v1/order", orderParams, &result)
	return result, err
}

func (client *RestClient) PlaceOrder3(orderParams OrderRequest3) (OrderResponse3, error) {
	var result OrderResponse3
	_, err := client.Request("POST", "/v3/order", orderParams, &result)
	return result, err
}

// ExecuteQuote executes the quote.
//         :param fx_quote_id: (str) the quote id received via get_quote
//         :param side: (str) must be either buy or sell
//             Example:
//                 {
//                     'status': 'success',
//                     'buy_price': 294.0,
//                     'error': None,
//                     'fx_quote_id': 'fad0ac687b1e439a92a0bafd92441e48',
//                     'is_filled': True,
//                     'price_executed': 294.0,
//                     'quantity_requested': {'token': 'ETH', 'value': '0.10000'},
//                     'sell_price': 293.94,
//                     'side_executed': 'buy',
//                     'side_requested': 'two_way',
//                     't_execute': '2019-07-03T21:45:10.358335+00:00',
//                     't_expiry': '2019-07-03T21:45:17.198692+00:00',
//                     't_quote': '2019-07-03T21:45:07.198688+00:00',
//                     'token_pair': {'base_token': 'ETH', 'quote_token': 'USD'}
//                 }
func (client *RestClient) ExecuteQuote(quoteParams QuoteExecutionRequest) (QuoteResponse, error) {
	var result QuoteResponse
	_, err := client.Request("POST", "/v1/quotes/execute", quoteParams, &result)
	return result, err
}

// GetQuoteStatus checks the status of a quote already requested.
//         :param fx_quote_id: (str) the quote id received via get_quote
//             Example:
//                 {
//                   "status": "success",
//                   "fx_quote_id": "00c884b056f949338788dfb59e495377",
//                   "buy_price": 12650,
//                   "sell_price": null,
//                   "platform": "api",
//                   "token_pair": {
//                     "base_token": "BTC",
//                     "quote_token": "USD"
//                   },
//                   "quantity_requested": {
//                     "token": "BTC",
//                     "value": "10"
//                   },
//                   "side_requested": "buy",
//                   "t_quote": "2019-06-27T11:59:21.875725+00:00",
//                   "t_expiry": "2019-06-27T11:59:22.875725+00:00",
//                   "is_filled": false,
//                   "side_executed": null,
//                   "price_executed": null,
//                   "t_execute": null,
//                   "trader_email": "trader1@company.com"
//                 }
func (client *RestClient) GetQuoteStatus(fxQuoteID string) (QuoteResponse, error) {
	var result QuoteResponse
	endPoint := fmt.Sprintf("/v1/quotes/%s", fxQuoteID)
	_, err := client.Request("GET", endPoint, nil, &result)
	return result, err
}

// GetExecutedQuotes gets a historical record of executed quotes in the time range.
//         :param t_start: (str) time in ISO8601 format (e.g. '2019-07-02T22:06:24.342342+00:00')
//         :param t_end: (str) time in ISO8601 format (e.g. '2019-07-03T22:06:24.234213+00:00'
//         :param platform: possible values -> ('browser', 'api', 'margin')
//             Example:
//                 [{'buy_price': 293.1, 'error': None, 'fx_quote_id': 'e2e1758f1a094a2a85825b592e9fc0d9',
//                 'is_filled': True, 'price_executed': 293.1, 'platform': 'browser', 'quantity_requested': {'token': 'ETH', 'value': '0.10000'},
//                 'sell_price': 293.03, 'side_executed': 'buy', 'side_requested': 'two_way', 'status': 'success',
//                 't_execute': '2019-07-03T14:02:56.539710+00:00', 't_expiry': '2019-07-03T14:03:02.038093+00:00',
//                 't_quote': '2019-07-03T14:02:52.038087+00:00',
//                 'token_pair': {'base_token': 'ETH', 'quote_token': 'USD'}, 'trader_email': 'trader1@company.com'},
//                 {'buy_price': 293.1, 'error': None, 'fx_quote_id': 'fc17a0d884444a0db5a7d9568c6c3f70',
//                 'is_filled': True, 'price_executed': 293.03, 'platform': 'api', 'quantity_requested': {'token': 'ETH', 'value': '0.10000'},
//                 'sell_price': 293.03, 'side_executed': 'sell', 'side_requested': 'two_way', 'status': 'success',
//                 't_execute': '2019-07-03T14:02:46.480337+00:00', 't_expiry': '2019-07-03T14:02:50.454222+00:00',
//                 't_quote': '2019-07-03T14:02:40.454217+00:00', 'token_pair': {'base_token': 'ETH', 'quote_token': 'USD'},
//                 'trader_email': 'trader2@company.com'}]
func (client *RestClient) GetExecutedQuotes(tStart time.Time, tEnd time.Time) ([]QuoteResponse, error) {
	var result []QuoteResponse
	requestParams := map[string]string{"t_start": tStart.Format(time.RFC3339), "t_end": tEnd.Format(time.RFC3339), "platform": "api"}
	_, err := client.Request("GET", "/v1/quotes", requestParams, &result)
	return result, err
}

func (client *RestClient) GetExecutedQuotesAll(tStart time.Time, tEnd time.Time) ([]QuoteResponse, error) {
	var result []QuoteResponse
	// NS: we don't specify the platform so we can get all trades
	requestParams := map[string]string{"t_start": tStart.Format(time.RFC3339), "t_end": tEnd.Format(time.RFC3339)}
	_, err := client.Request("GET", "/v1/quotes", requestParams, &result)
	return result, err
}

// GetBalances gets account balances.
//         :param platform: possible values -> ('browser', 'api', 'margin')
//             Example:
//                 [
//                     {'balance': 0.0, 'token': 'BTC', 'platform': 'browser'},
//                     {'balance': -1.3772005993291505, 'token': 'ETH', 'platform': 'api'},
//                     {'balance': 187.624207, 'token': 'USD', 'platform': 'api'}
//                 ]
func (client *RestClient) GetBalances() ([]Balance, error) {
	var result []Balance
	requestParams := map[string]string{"platform": "api"}
	_, err := client.Request("GET", "/v1/balances", requestParams, &result)
	return result, err
}

// GetTransfers gets a historical record of deposits/withdrawals between the given time range.
//         :param t_start: (str) time in ISO8601 format (e.g. '2019-07-02T22:06:24.342342+00:00')
//         :param t_end: (str) time in ISO8601 format (e.g. '2019-07-03T22:06:24.234213+00:00'
//         :param platform: possible values -> ('browser', 'api', 'margin')
//             Example:
//                 [
//                   {
//                     "type": "deposit",
//                     "platform": "api",
//                     "token": "BTC",
//                     "quantity": 1.0,
//                     "t_create": "2019-06-20T01:01:01+00:00"
//                   },
//                   {
//                     "type": "withdrawal",
//                     "platform": "midas",
//                     "token": "BTC",
//                     "quantity": 1.0,
//                     "t_create": "2019-06-22T01:01:01+00:00"
//                   }
//                 ]
func (client *RestClient) GetTransfers(tStart time.Time, tEnd time.Time) ([]Transfer, error) {
	var result []Transfer
	requestParams := map[string]string{"t_start": tStart.Format(time.RFC3339), "t_end": tEnd.Format(time.RFC3339)}
	_, err := client.Request("GET", "/v1/transfers", requestParams, &result)
	return result, err
}

func (client *RestClient) GetTradeVolume(tStart time.Time, tEnd time.Time) (TradeVolume, error) {
	var result TradeVolume
	requestParams := map[string]string{"t_start": tStart.Format(time.RFC3339), "t_end": tEnd.Format(time.RFC3339), "platform": "api"}
	_, err := client.Request("GET", "/v1/get_trade_volume", requestParams, &result)
	return result, err
}

func (client *RestClient) GetTradeLimits(platform string) (TradeLimits, error) {
	var result TradeLimits
	endPoint := fmt.Sprintf("/v1/get_trade_limits/%s", platform)
	_, err := client.Request("GET", endPoint, nil, &result)
	return result, err
}

func (client *RestClient) GetTradeSizes() ([]TradeSize, error) {
	var result []TradeSize
	_, err := client.Request("GET", "/v1/trade_sizes", nil, &result)
	return result, err
}

func (client *RestClient) GetTotalBalances() ([]TotalBalance, error) {
	var result []TotalBalance
	_, err := client.Request("GET", "/v1/balances/total", nil, &result)
	return result, err
}
