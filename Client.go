package cryptoprice

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

const (
	cryptoCompareURL = "https://min-api.cryptocompare.com/data"
	nowURI           = "/price"
	minuteURI        = "/histominute"
	hourURI          = "/histohour"
)

type Client struct {
	URL           string
	FromSymbol    string
	ToSymbol      string
	TryConversion bool
	Exchange      string
	ExtraParams   string
	http.Client
}

func NewClient(fromSymbol, toSymbol string) *Client {
	return &Client{
		FromSymbol:    fromSymbol,
		ToSymbol:      toSymbol,
		TryConversion: true,
		URL:           cryptoCompareURL,
		ExtraParams:   "golang pkg - github.com/canonical-ledgers/cryptoprice",
	}
}

func (c Client) GetPrice(t time.Time) (float64, error) {
	since := time.Since(t)
	var values url.Values
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

	URI := hourURI
	response := interface{}(historicalResponseT{})
	if since < 1*time.Minute {
		values.Add("tsyms", c.ToSymbol)
		URI = nowURI
		response = nowResponseT{}
	} else {
		values.Add("tsym", c.ToSymbol)
		values.Add("toTs", c.ToSymbol)
		values.Add("limit", fmt.Sprintf("%v", 1))
	}
	if since < 7*24*time.Hour {
		URI = minuteURI
	}
	req, err := http.NewRequest("GET", c.URL+URI+values.Encode(), nil)
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
	case nowResponseT:
		if v, ok := r[c.ToSymbol]; ok {
			if v, ok := v.(float64); ok {
				return v, nil
			}
		}
		return 0, fmt.Errorf("Unknown response type: %+v", r)
	case historicalResponseT:
		if r.Response != "Success" {
			return 0, fmt.Errorf("Response: %#v, Message: %#v",
				r.Response, r.Message)
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
