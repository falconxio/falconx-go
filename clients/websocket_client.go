package clients

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"net/http"
	"time"
)

const (
	pingInterval            = 20 * time.Second
	webSocketSecureProtocol = "wss://"
	socketioUrl             = "/socket.io/?EIO=3&transport=websocket"
)

type SocketClientConfig struct {
	Host       string
	Secret     string
	APIKey     string
	Passphrase string
}

type SocketClient struct {
	Config     SocketClientConfig
	Namespace  string
	Transport  *transport.WebsocketTransport
	Connection *gosocketio.Client
}

func ComputeHmac256(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func NewSocketClient(config SocketClientConfig, namespace string) *SocketClient {
	wst := transport.GetDefaultWebsocketTransport()
	wst.PingInterval = pingInterval
	wst.RequestHeader = make(http.Header)

	client := SocketClient{
		Config:     config,
		Namespace:  namespace,
		Transport:  wst,
		Connection: nil,
	}

	return &client
}

func (client *SocketClient) Connect() error {
	err := client.AddAuth()
	if err != nil {
		log.Fatal("Error creating authentication parameters.")
		return err
	}
	falconxWsUrl := webSocketSecureProtocol + client.Config.Host + socketioUrl
	client.Connection, err = gosocketio.Dial(falconxWsUrl, client.Transport)
	return err
}

func (client *SocketClient) AddAuth() error {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	method := "GET"
	path := "/socket.io/"
	message := timestamp + method + path
	hmac_key, err := base64.StdEncoding.DecodeString(client.Config.Secret)
	if err != nil {
		return err
	}
	signature := ComputeHmac256(message, hmac_key)

	client.Transport.RequestHeader.Add("FX-ACCESS-SIGN", signature)
	client.Transport.RequestHeader.Add("FX-ACCESS-TIMESTAMP", timestamp)
	client.Transport.RequestHeader.Add("FX-ACCESS-KEY", client.Config.APIKey)
	client.Transport.RequestHeader.Add("FX-ACCESS-PASSPHRASE", client.Config.Passphrase)
	client.Transport.RequestHeader.Add("Content-Type", "application/json")
	return nil
}
