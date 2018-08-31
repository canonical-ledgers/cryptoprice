// github.com/canonical-ledgers/cryptoprice v1.0.0
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package cryptoprice_test

import (
	"fmt"
	"time"

	"github.com/canonical-ledgers/cryptoprice"
)

func ExampleClient_GetPrice() {
	client := cryptoprice.NewClient("BTC", "USD")
	client.Timeout = 5 * time.Second
	p, err := client.GetPrice(time.Now())
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Latest BTC price: $%v USD\n", p)
}
