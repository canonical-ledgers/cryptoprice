// github.com/canonical-ledgers/cryptoprice v1.0.0
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package cryptoprice

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

// Default URIs and parameters
const (
	CryptoCompareURL = "https://min-api.cryptocompare.com/data"
	NowURI           = "/price"
	MinuteURI        = "/histominute"
	HourURI          = "/histohour"
	ExtraParams      = "golang pkg - github.com/canonical-ledgers/cryptoprice"
)

// Client stores request parameters and provides methods for querying the
// CryptoCompare REST API. You may set any additional http.Client settings
// directly on this type.
type Client struct {
	URL           string // URL to send requests to
	FromSymbol    string // Cryptocurrency symbol of interest
	ToSymbol      string // Currency symbol to convert into
	TryConversion bool   // Set to false to only return data if a direct pair is available
	Exchange      string // Exchange to use for data (default: "CCCAGG" aggregated average)
	ExtraParams   string // Name of your application, defaults to const ExtraParams
	http.Client
}

// NewClient returns a pointer to a new Client with the given FromSymbol and
// ToSymbol set, as well as the correct CryptoCompare API endpoint URL,
// TryConversion set to true, and the default ExtraParams set.
func NewClient(fromSymbol, toSymbol string) *Client {
	return &Client{
		FromSymbol:    fromSymbol,
		ToSymbol:      toSymbol,
		TryConversion: true,
		URL:           CryptoCompareURL,
		ExtraParams:   ExtraParams,
	}
}

// GetPrice returns the most accurate price available for the given time t. If
// the requested time is within the past minute, the most recent price data is
// used. If the requested time is within the past 7 days, the simple average of
// the high and low prices for the minute that is closest to the given time is
// used. If the request time is any further in the past, the simple average of
// the high and low prices for the hour that is closest to the given price is
// used.
func (c Client) GetPrice(t time.Time) (float64, error) {
	if len(c.FromSymbol) == 0 || len(c.ToSymbol) == 0 {
		return 0, fmt.Errorf("FromSymbol and ToSymbol not specified")
	}
	values := make(url.Values)
	values.Add("fsym", c.FromSymbol)
	if len(c.ExtraParams) > 0 {
		values.Add("extraParams", c.ExtraParams)
	}
	if len(c.Exchange) > 0 {
		values.Add("e", c.Exchange)
	}
	if !c.TryConversion {
		values.Add("tryConversion", fmt.Sprintf("%v", c.TryConversion))
	}

	var response interface{}
	response = &historicalResponseT{}

	URI := HourURI
	since := time.Since(t)
	roundUp := time.Hour
	if since < 7*24*time.Hour {
		URI = MinuteURI
		roundUp = time.Minute
	}
	if since < 1*time.Minute {
		values.Add("tsyms", c.ToSymbol)
		URI = NowURI
		response = make(map[string]interface{})
	} else {
		values.Add("tsym", c.ToSymbol)
		values.Add("toTs", fmt.Sprintf("%v",
			t.Truncate(roundUp).Add(roundUp).Unix()))
		values.Add("limit", fmt.Sprintf("%v", 1))
	}

	req, err := http.NewRequest("GET", c.URL+URI+"?"+values.Encode(), nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&response); err != nil {
		return 0, err
	}

	switch r := response.(type) {
	case map[string]interface{}:
		if v, ok := r[c.ToSymbol]; ok {
			if v, ok := v.(float64); ok {
				return v, nil
			}
		}
		return 0, fmt.Errorf("Unknown response type: %+v", r)
	case *historicalResponseT:
		if r.Response != "Success" {
			return 0, fmt.Errorf("Response: %#v, Message: %#v",
				r.Response, r.Message)
		}
		if len(r.Data) == 0 {
			return 0, fmt.Errorf("No data returned")
		}
		var pID int
		minDuration := time.Duration(math.MaxInt64)
		for i, p := range r.Data {
			duration := t.Sub(p.Time)
			if duration < minDuration {
				minDuration = duration
				pID = i
			}
		}
		price := r.Data[pID]
		return (price.High + price.Low) / 2, nil
	}
	panic("Unreachable code")
}

type nowResponseT map[string]interface{}
type historicalResponseT struct {
	Response string
	Message  string
	Data     []priceT
}

type priceT struct {
	Time time.Time
	High float64
	Low  float64
}
