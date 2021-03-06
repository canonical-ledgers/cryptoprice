// github.com/canonical-ledgers/cryptoprice
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package cryptoprice_test

import (
	"fmt"
	"time"

	"github.com/canonical-ledgers/cryptoprice/v2"
)

func ExampleClient_GetPriceAt() {
	client := cryptoprice.NewClient("BTC", "USD")
	client.Timeout = 5 * time.Second
	p, err := client.GetPriceAt(time.Now())
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Latest BTC price: $%v USD\n", p)
}
