/*
 * Copyright (C) 2019 The onyxchain Authors
 * This file is part of The onyxchain library.
 *
 * The onyxchain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The onyxchain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The onyxchain.  If not, see <http://www.gnu.org/licenses/>.
 */

package constants

import (
	"time"
)

// genesis constants
var (
	//TODO: modify this when on mainnet
	GENESIS_BLOCK_TIMESTAMP = uint32(time.Date(2018, time.June, 30, 0, 0, 0, 0, time.UTC).Unix())
)

// onx constants
const (
	ONX_NAME           = "ONYX Token"
	ONX_SYMBOL         = "ONYX"
	ONX_DECIMALS       = 1 * 8
	ONX_DECIMALS_RATIO = uint64(100000000) // 10 ** 8
	ONX_RATIO          = uint64(10 * ONX_DECIMALS_RATIO)
	ONX_TOTAL_SUPPLY   = uint64(1000000000) * ONX_RATIO
)

// oxg constants
const (
	OXG_NAME         = "OXG Token"
	OXG_SYMBOL       = "OXG"
	OXG_DECIMALS     = 9
	OXG_TOTAL_SUPPLY = uint64(1000000000000000000)
)

// onx/oxg unbound model constants
const UNBOUND_TIME_INTERVAL = uint32(31536000)

var UNBOUND_GENERATION_AMOUNT = [18]uint64{5, 4, 3, 3, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

// the end of unbound timestamp offset from genesis block's timestamp
var UNBOUND_DEADLINE = (func() uint32 {
	count := uint64(0)
	for _, m := range UNBOUND_GENERATION_AMOUNT {
		count += m
	}
	count *= uint64(UNBOUND_TIME_INTERVAL)

	numInterval := len(UNBOUND_GENERATION_AMOUNT)

	if UNBOUND_GENERATION_AMOUNT[numInterval-1] != 1 ||
		(count-uint64(UNBOUND_TIME_INTERVAL) > ONX_TOTAL_SUPPLY/ONX_RATIO && ONX_TOTAL_SUPPLY/ONX_RATIO > count) {
		panic("incompatible constants setting")
	}

	return UNBOUND_TIME_INTERVAL*uint32(numInterval) - uint32(count-uint64(ONX_TOTAL_SUPPLY)/ONX_RATIO)
})()

// multi-sig constants
const MULTI_SIG_MAX_PUBKEY_SIZE = 16

// transaction constants
const TX_MAX_SIG_SIZE = 16

// network magic number
const (
	NETWORK_MAGIC_MAINNET = 0x8c77ab60
	NETWORK_MAGIC_POLARIS = 0x2d8829df
)
