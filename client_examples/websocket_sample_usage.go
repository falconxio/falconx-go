package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	// "github.com/paxosglobal/golang-socketio"
	// "github.com/paxosglobal/golang-socketio/transport"
	"net/http"
	"time"

	"go.uber.org/multierr"
)

const (
	API_KEY            = ""
	SECRET_KEY         = ""
	PASSPHRASE         = ""
	falconxWsUrl       = "wss://ws-bikram.falconxdev.com/?transport=websocket&EIO=3"
	streamingNamespace = "/streaming"
	pingInterval       = 20 * time.Second
)

type Credentials struct {
	ApiKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
	PassPhrase string `json:"passphrase"`
}

type SubscriptionRequest struct {
	TokenPair       TokenPair `json:"token_pair"`
	Quantity        []string  `json:"quantity"`
	ClientRequestID string    `json:"client_request_id"`
}

type TokenPair struct {
	BaseToken  string `json:"base_token"`
	QuoteToken string `json:"quote_token"`
}

func AddFalconxAuthHeader(header http.Header, creds *Credentials, method string, path string, body string) error {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	message := timestamp + method + path + body
	hmac_key, err := base64.StdEncoding.DecodeString(creds.SecretKey)
	if err != nil {
		return err
	}
	signature := ComputeHmac256(message, hmac_key)

	header.Add("FX-ACCESS-SIGN", signature)
	header.Add("FX-ACCESS-TIMESTAMP", timestamp)
	header.Add("FX-ACCESS-KEY", creds.ApiKey)
	header.Add("FX-ACCESS-PASSPHRASE", creds.PassPhrase)
	header.Add("Content-Type", "application/json")
	return nil
}

func ComputeHmac256(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	creds := &Credentials{API_KEY, SECRET_KEY, PASSPHRASE}
	wst := transport.GetDefaultWebsocketTransport()
	wst.PingInterval = pingInterval
	wst.RequestHeader = make(http.Header)
	path := "/socket.io/"
	_ = AddFalconxAuthHeader(wst.RequestHeader, creds, "GET", path, "")
	client, err := gosocketio.Dial(falconxWsUrl, wst)
	if err != nil {
		panic(err)
	}

	client.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		fmt.Println("falconx socketio server connected")
	})

	client.On("response", func(c *gosocketio.Channel, msg interface{}) {
		fmt.Println("falconx socketio server response.", "message:", msg)
	})

	client.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel, reasons gosocketio.ConnectionErrors) {
		errors := multierr.Combine(reasons.Errors...)
		fmt.Println("falconx socketio server disconnected.", "error:", errors)
	})

	client.ConnectNamespace(streamingNamespace)

	client.On("stream", func(c *gosocketio.Channel, stream interface{}) {
		fmt.Println("price update received", "stream", stream)
	})

	client.Emit("subscribe", streamingNamespace, &SubscriptionRequest{
		TokenPair: TokenPair{
			BaseToken:  "BTC",
			QuoteToken: "USD",
		},
		Quantity:        []string{"0.000009", "0.00009", "0.0009", "0.009", "0.09", "0.9"},
		ClientRequestID: "1",
	})
	//... subscribe to other markets

	time.Sleep(10 * time.Minute)
}
