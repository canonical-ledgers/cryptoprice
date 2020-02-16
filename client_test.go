// github.com/canonical-ledgers/cryptoprice
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package cryptoprice

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPriceAt(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	reqT = t
	testServer := httptest.NewServer(handler)

	client := NewClient("", "USD")
	client.URL = testServer.URL
	_, err := client.GetPriceAt(time.Now())
	assert.Error(err, "Empty FromSymbol")

	client.Currency = "BTC"
	client.Units = ""
	_, err = client.GetPriceAt(time.Now())
	assert.Error(err, "Empty ToSymbol")

	client.Units = "USD"
	msg := "GetPriceAt(time.Now())"
	p, err := client.GetPriceAt(time.Now())
	require.NoError(err, msg)
	assert.Equal(5.5, p, msg)
	assert.Equal(reqURL.Path, NowURI, msg)

	msg = "GetPriceAt(time.Now().Add(-2 * time.Minute))"
	p, err = client.GetPriceAt(time.Now().Add(-2 * time.Minute))
	require.NoError(err, msg)
	assert.Equal(5.5, p, msg)
	assert.Equal(reqURL.Path, MinuteURI, msg)

	msg = "GetPriceAt(time.Now().Add(-8 * 24 * time.Hour))"
	p, err = client.GetPriceAt(time.Now().Add(-8 * 24 * time.Hour))
	require.NoError(err, msg)
	assert.Equal(5.5, p, msg)
	assert.Equal(reqURL.Path, HourURI, msg)
}

var reqURL *url.URL
var reqForm url.Values
var reqT *testing.T
var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	require := require.New(reqT)
	require.NoError(r.ParseForm(), "handler: http.Request.ParseForm()")
	reqURL = r.URL
	reqForm = r.Form
	require.NotEmpty(reqForm, "handler: http.Request.Form")
	var resp interface{}
	switch r.URL.Path {
	case NowURI:
		resp = map[string]float64{r.FormValue("tsyms"): 5.5}
	case MinuteURI:
		fallthrough
	case HourURI:
		unixSec, err := strconv.ParseInt(r.FormValue("toTs"), 10, 64)
		require.NoErrorf(err, "strconv.Atoi(toTs), toTs: %#v", r.FormValue("toTs"))
		ts := time.Unix(unixSec, 0).Truncate(time.Minute)
		resp = historicalResponseT{
			Response: "Success",
			Data: []priceT{{
				Time: timeT(ts.Add(-1 * time.Minute)),
				High: 5,
				Low:  4,
			}, {
				Time: timeT(ts),
				High: 6,
				Low:  5,
			}},
		}
	default:
		require.FailNowf("handler: Unrecognized http.Request.Path",
			"%#v", r.URL.Path)
	}
	e := json.NewEncoder(w)
	require.NoError(e.Encode(resp), "json.Encoder.Encode", "%+v", resp)
})
