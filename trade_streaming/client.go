package trade_streaming

//go:generate mockgen -destination client_mock.go -package trade_streaming github.com/fernandovaladao/vwap/trade_streaming Client
type Client interface {
	ReadValue() (*Trade, error)
}

// the first attribute is used to validate if this is a match or an error message.
// the second and third attributes are used to log error messages.
type Trade struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Pair    string `json:"product_id"`
	Price   string `json:"price"`
}
