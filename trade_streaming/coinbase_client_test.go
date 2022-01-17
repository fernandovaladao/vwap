package trade_streaming

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	testCases := []struct {
		name string
		mockConn func(ctrl *gomock.Controller) Conn
		doAsserts func(t *testing.T, trade *Trade, err error)
	}{
		{
			name: "Happy path: no errors returned",
			mockConn: func(ctrl *gomock.Controller) Conn {
				conn := NewMockConn(ctrl)
				conn.EXPECT().ReadJSON(gomock.Any())
				return conn
			},
			doAsserts: func(t *testing.T, trade *Trade, err error) {
				assert.NotNil(t, trade)
				assert.Nil(t, err)
			},
		},
		{
			name: "Error returned when connection reads",
			mockConn: func(ctrl *gomock.Controller) Conn {
				conn := NewMockConn(ctrl)
				conn.EXPECT().ReadJSON(gomock.Any()).Return(fmt.Errorf("error"))
				return conn
			},
			doAsserts: func(t *testing.T, trade *Trade, err error) {
				assert.Nil(t, trade)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			conn := tt.mockConn(ctrl)
			client := CoinbaseClient{
				conn: conn,
			}

			trade, err := client.ReadValue()

			tt.doAsserts(t, trade, err)
		})
	}
}