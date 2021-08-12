package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/falconxio/falconx-go/clients"
)

func main() {
	apiKey := "XXX"
	passPhrase := "XXX"
	secret := "XXX"

	baseUrl := "https://qa.falconxdev.com"

	client := clients.NewRestClient(clients.RestClientConfig{
		BaseURL:    baseUrl,
		APIKey:     apiKey,
		Passphrase: passPhrase,
		Secret:     secret,
	})

	tokenPair := clients.TokenPair{BaseToken: "BTC", QuoteToken: "USD"}
	quantity := clients.Quantity{Token: "BTC", Value: 0.001}
	quoteParams := clients.QuoteRequest{TokenPair: tokenPair, Quantity: quantity, Side: "buy", ClientOrderId: "343434343er4"}
	quoteResponse, err := client.GetQuote(quoteParams)
	quoteResponseJson, _ := json.Marshal(quoteResponse)

	fmt.Printf("\n\n Quote Resp: \n%s, err: %+v", quoteResponseJson, err)

	fxQuoteId := quoteResponse.FxQuoteId
	limit_price := quoteResponse.BuyPrice + 5
	side := quoteResponse.SideRequested

	// Quote Not Executed Here
	quoteStatus, err := client.GetQuoteStatus(fxQuoteId)
	quoteStatusJson, _ := json.Marshal(quoteStatus)
	fmt.Printf("\n\n %s Status: \n %s \n, err: %+v", fxQuoteId, quoteStatusJson, err)

	quoteExRequest := clients.QuoteExecutionRequest{FxQuoteId: fxQuoteId, Side: side}
	executionResponse, err := client.ExecuteQuote(quoteExRequest)
	executionResponseJson, _ := json.Marshal(executionResponse)
	fmt.Printf("\n\nQuote Execution Response: \n%s, \n error: %+v", executionResponseJson, err)

	// Quote Executed Here
	quoteStatus, err = client.GetQuoteStatus(fxQuoteId)
	quoteStatusJson, _ = json.Marshal(quoteStatus)
	fmt.Printf("\n\n %s Status: \n %s \n, err: %+v", fxQuoteId, quoteStatusJson, err)

	orderParams := clients.OrderRequest{TokenPair: tokenPair, Quantity: quantity, Side: "buy", OrderType: "limit", TimeInForce: "fok", LimitPrice: limit_price, SlippageBps: 5, ClientOrderId: "Bazinga"}
	orderResponse, err := client.PlaceOrder(orderParams)
	orderResponseJson, _ := json.Marshal(orderResponse)
	fmt.Printf("\n\nOrder Execution Response : \n%s\n error: %+v", orderResponseJson, err)

	balances, err := client.GetBalances()
	balancesJson, _ := json.Marshal(balances)
	fmt.Printf("\n\nBalances: \n%s, \nerr: %+v", balancesJson, err)

	validTokenPairs, err := client.GetTradingPairs()
	validTokenPairsJson, _ := json.Marshal(validTokenPairs)
	fmt.Printf("\n\nTrading Pairs : \n%s\n error: %+v", validTokenPairsJson, err)

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)

	executions, err := client.GetExecutedQuotes(startTime, endTime)
	executionsJson, _ := json.Marshal(executions)
	fmt.Printf("\n\nexecutions: \n%s, err: %+v", executionsJson, err)

	totalBalance, err := client.GetTotalBalances()
	fmt.Printf("\n\nTotal Balances: \n%+v, err: %+v", totalBalance, err)

	tradeSizes, err := client.GetTradeSizes()
	tradeSizesJson, _ := json.Marshal(tradeSizes[0])
	fmt.Printf("\n\n Trade Sizes: \n%s, \n err: %+v", tradeSizesJson, err)

	tradeLimits, err := client.GetTradeLimits("api")
	fmt.Printf("\n\n Trade Limits: \n%+v, \n err: %+v", tradeLimits, err)

	tradeVolumes, err := client.GetTradeVolume(startTime, endTime)
	fmt.Printf("\n\n Trade Volumes: \n%+v, \n err: %+v", tradeVolumes, err)

	transfers, err := client.GetTransfers(startTime, endTime)
	fmt.Printf("\n\n Transfers: \n%+v, \n err: %+v", transfers, err)
}
