package client_examples

import (
	"fmt"
	"log"
	"time"

	"github.com/falconxio/falconx-go/clients"
	gosocketio "github.com/graarh/golang-socketio"
	"go.uber.org/multierr"
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

func RunWebSocketExamples(apiKey string, secret string, passphrase string, host string) {

	if len(host) == 0 {
		host = "ws.falconx.io"
	}

	streamingNamespace := "/streaming"

	client := clients.NewSocketClient(
		clients.SocketClientConfig{
			Host:       host,
			Secret:     secret,
			APIKey:     apiKey,
			Passphrase: passphrase,
		},
		streamingNamespace,
	)

	err := client.Connect()
	if err != nil {
		log.Println(err)
		log.Fatal("Error Connecting to Falconx!")
		panic(err)
	}

	client.Connection.On(OnConnection, func(channel *gosocketio.Channel) {
		fmt.Println("Connection established")
	})

	client.Connection.On(OnDisconnection, func(channel *gosocketio.Channel, reasons gosocketio.ConnectionErrors) {
		errors := multierr.Combine(reasons.Errors...)
		fmt.Println("Disconnected from FalconX", "error:", errors)
	})

	client.Connection.ConnectNamespace(streamingNamespace)

	client.Connection.On(OnResponse, func(channel *gosocketio.Channel, msg interface{}) {
		fmt.Println("Response received from Socket: ", "message:", msg)
	})

	client.Connection.On(OnStream, func(channel *gosocketio.Channel, stream interface{}) {
		fmt.Println("Price change tick received", "stream", stream)
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
			BaseToken:  "ETH",
			QuoteToken: "USD",
		},
		Quantity:        []float64{0.001, 0.01, 0.1},
		ClientRequestID: "5c5325e3-ee42-76fa-932c-64dce446d8be",
	})

	wait:= make(chan bool, 1)
	<-wait
}
