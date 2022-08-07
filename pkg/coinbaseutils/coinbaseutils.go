package coinbaseutils

/*
Package contaning functions fetching from the coinbase api.
*/

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/parthrs/btctracker/pkg/httputils"
)

// GetBtcUsdPrice returns the spot price for BTC-USD as a string and an
// error.
func GetBtcUsdPrice() (price float64, err error) {
	resp, err := httputils.HttpGet("https://api.coinbase.com/v2/prices/BTC-USD/spot")
	if err != nil {
		return
	}
	var respMap map[string]map[string]string
	json.Unmarshal(resp, &respMap)
	p, ok := respMap["data"]["amount"]
	if !ok {
		err = fmt.Errorf("empty response from coinbase api")
		return
	}

	price, err = strconv.ParseFloat(p, 64)
	return
}
