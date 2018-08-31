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
