package trade_streaming

import (
	"github.com/gorilla/websocket"
)

const websocketFeed = "wss://ws-feed-public.sandbox.exchange.coinbase.com"

// since Gorilla's websocket.Conn is not an interface, we can't use mockGen to mock it. Hence, we are wrapping it
// around this interface to bypass this limitation.
//go:generate mockgen -destination conn_mock.go -package trade_streaming github.com/fernandovaladao/vwap/trade_streaming Conn
type Conn interface {
	ReadJSON(v interface{}) error
}

type coinbaseConn struct {
	conn *websocket.Conn
}

func (cc coinbaseConn) ReadJSON(v interface{}) error {
	return cc.conn.ReadJSON(v)
}

type CoinbaseClient struct {
	tradingPairs []string
	conn         Conn
}

type subscriptionMsg struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

func NewCoinbaseClient(tradingPairs []string) (*CoinbaseClient, error) {
	// first creates connection to the Coinbase Websocket Feed
	if conn, _, err := websocket.DefaultDialer.Dial(websocketFeed, nil); err != nil {
		return nil, err
	} else {
		// next sends subscription message with trading pairs and matches channel
		subscriptionMsg := newSubscriptionMsg(tradingPairs)
		if err = conn.WriteJSON(subscriptionMsg); err != nil {
			return nil, err
		}

		return &CoinbaseClient{
			tradingPairs: tradingPairs,
			conn: coinbaseConn{
				conn: conn,
			},
		}, nil
	}
}

func (cc *CoinbaseClient) ReadValue() (*Trade, error) {
	trade := &Trade{}
	err := cc.conn.ReadJSON(trade)
	if err != nil {
		return nil, err
	}
	return trade, nil
}

func newSubscriptionMsg(tradingPairs []string) subscriptionMsg {
	return subscriptionMsg{
		Type:       "subscribe",
		ProductIds: tradingPairs,
		Channels:   []string{"matches"},
	}
}
