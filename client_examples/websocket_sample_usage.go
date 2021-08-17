package main

import (
	"fmt"
	"github.com/falconxio/falconx-go/clients"
	gosocketio "github.com/graarh/golang-socketio"
	"go.uber.org/multierr"
	"log"
	"time"
)

const (
	OnConnection    = "connection"
	OnDisconnection = "disconnection"
	OnResponse      = "response"
	OnStream        = "stream"
	OnError         = "error"

	MaxConnectionRequest  = "GET_MAX_CONNECTIONS"
	AllowedMarketsRequest = "GET_ALLOWED_MARKETS"
	MaxLevelsRequest      = "GET_MAX_LEVELS"
)

func main() {

	host := "ws.falconx.io"
	API_KEY := ""
	SECRET_KEY := ""
	PASSPHRASE := ""
	streamingNamespace := "/streaming"

	client := clients.NewSocketClient(
		clients.SocketClientConfig{
			Host:       host,
			Secret:     SECRET_KEY,
			APIKey:     API_KEY,
			Passphrase: PASSPHRASE,
		},
		streamingNamespace,
	)

	err := client.Connect()
	if err != nil {
		log.Fatal("Error Connecting to Falconx!")
		panic(err)
	}

	client.Connection.On(OnConnection, func(channel *gosocketio.Channel) {
		fmt.Println("falconx socketio server connected")
	})

	client.Connection.On(OnDisconnection, func(channel *gosocketio.Channel, reasons gosocketio.ConnectionErrors) {
		errors := multierr.Combine(reasons.Errors...)
		fmt.Println("falconx socketio server disconnected.", "error:", errors)
	})

	client.Connection.ConnectNamespace(streamingNamespace)

	client.Connection.On(OnResponse, func(channel *gosocketio.Channel, msg interface{}) {
		fmt.Println("falconx socketio server response.", "message:", msg)
	})

	client.Connection.On(OnStream, func(channel *gosocketio.Channel, stream interface{}) {
		fmt.Println("price update received", "stream", stream)
	})

	// User Config Requests
	client.Connection.Emit("request", streamingNamespace, &clients.UserConfigRequest{
		MessageType:     AllowedMarketsRequest,
		ClientRequestID: "5c5325e3-ee42-76fa-932c-64dce446d8be",
	})

	client.Connection.Emit("request", streamingNamespace, &clients.UserConfigRequest{
		MessageType:     MaxConnectionRequest,
		ClientRequestID: "5c5325e3-ee42-76fa-932c-64dce446d8be",
	})

	client.Connection.Emit("request", streamingNamespace, &clients.UserConfigRequest{
		MessageType:     MaxLevelsRequest,
		ClientRequestID: "5c5325e3-ee42-76fa-932c-64dce446d8be",
	})

	client.Connection.Emit("subscribe", streamingNamespace, &clients.SubscriptionRequest{
		TokenPair: clients.TokenPair{
			BaseToken:  "BTC",
			QuoteToken: "USD",
		},
		Quantity:        []string{"0.000009", "0.00009", "0.0009", "0.009", "0.09", "0.9"},
		ClientRequestID: "1",
	})

	time.Sleep(10 * time.Minute)
}
