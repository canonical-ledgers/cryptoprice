// github.com/canonical-ledgers/cryptoprice
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
	CryptoCompareURL   = "https://min-api.cryptocompare.com/data"
	NowURI             = "/price"
	MinuteURI          = "/histominute"
	HourURI            = "/histohour"
	DefaultExtraParams = "golang pkg - github.com/canonical-ledgers/cryptoprice"
)

// Client stores request parameters and provides methods for querying the
// CryptoCompare REST API. You may set any additional http.Client settings
// directly on this type.
type Client struct {
	URL string

	// Get price of Currency in Units.
	Currency string
	Units    string

	// Only return data for direct trading pairs. Do not perform conversions.
	DirectPairOnly bool

	// Exchange to use for data (default: "CCCAGG" aggregated average)
	Exchange string

	// ExtraParams should be the name of your application, defaults to
	// const DefaultExtraParams
	ExtraParams string
	http.Client
}

// NewClient returns a pointer to a new Client with the given Currency and
// Units set, as well as the correct CryptoCompare API endpoint URL.
func NewClient(currency, units string) *Client {
	return &Client{
		Currency:    currency,
		Units:       units,
		URL:         CryptoCompareURL,
		ExtraParams: DefaultExtraParams,
	}
}

func (c *Client) GetPriceNow() (float64, error) {
	return c.GetPriceAt(time.Now())
}

// GetPriceAt returns the most accurate price available for the given time t.
// If the requested time is within the past minute, the most recent price data
// is used. If the requested time is within the past 7 days, the simple average
// of the high and low prices for the minute that is closest to the given time
// is used. If the request time is any further in the past, the simple average
// of the high and low prices for the hour that is closest to the given price
// is used.
func (c *Client) GetPriceAt(t time.Time) (float64, error) {
	if len(c.Currency) == 0 || len(c.Units) == 0 {
		return 0, fmt.Errorf("Currency and Units not specified")
	}
	values := make(url.Values)
	values.Add("fsym", c.Currency)
	if len(c.ExtraParams) > 0 {
		values.Add("extraParams", c.ExtraParams)
	}
	if len(c.Exchange) > 0 {
		values.Add("e", c.Exchange)
	}
	if c.DirectPairOnly {
		values.Add("tryConversion", "false")
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
		values.Add("tsyms", c.Units)
		URI = NowURI
		response = make(map[string]interface{})
	} else {
		values.Add("tsym", c.Units)
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
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&response); err != nil {
		return 0, err
	}
	switch r := response.(type) {
	case map[string]interface{}:
		if v, ok := r[c.Units]; ok {
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
			duration := t.Sub(time.Time(p.Time))
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
	Time timeT
	High float64
	Low  float64
}

type timeT time.Time

func (t *timeT) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	*t = timeT(time.Unix(timestamp, 0))
	return nil
}

func (t timeT) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", time.Time(t).Unix())), nil
}
