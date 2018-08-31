// github.com/canonical-ledgers/cryptoprice v1.0.0
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

// Package cryptoprice provides a simple Client for querying the most accurate
// crypto currency price available from the CryptoCompare REST API.
//
// This is not intended to be a complete implementation of the full
// CryptoCompare REST API. It just accesses the CryptoCompare REST API
// endpoints required to get the most accurate price available for a given
// time.
//
// The CryptoCompare REST API documentation can be found here:
// https://min-api.cryptocompare.com/
package cryptoprice
