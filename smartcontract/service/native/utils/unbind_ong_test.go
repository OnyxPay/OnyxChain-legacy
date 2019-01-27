/*
 * Copyright (C) 2018 The OnyxChain Authors
 * This file is part of The OnyxChain library.
 *
 * The OnyxChain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OnyxChain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The OnyxChain.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"math/rand"
	"testing"

	"github.com/OnyxPay/OnyxChain-legacy/common/constants"
	"github.com/stretchr/testify/assert"
)

func TestCalcUnbindOxg(t *testing.T) {
	assert.Equal(t, CalcUnbindOxg(1, 0, 1), uint64(GENERATION_AMOUNT[0]))
	assert.Equal(t, CalcUnbindOxg(1, 0, TIME_INTERVAL), GENERATION_AMOUNT[0]*uint64(TIME_INTERVAL))
	assert.Equal(t, CalcUnbindOxg(1, 0, TIME_INTERVAL+1),
		GENERATION_AMOUNT[1]+GENERATION_AMOUNT[0]*uint64(TIME_INTERVAL))
}

// test identity: unbound[t1, t3) = unbound[t1, t2) + unbound[t2, t3)
func TestCumulative(t *testing.T) {
	N := 10000
	for i := 0; i < N; i++ {
		tstart := rand.Uint32()
		tend := tstart + rand.Uint32()
		tmid := uint32((uint64(tstart) + uint64(tend)) / 2)

		total := CalcUnbindOxg(1, tstart, tend)
		total2 := CalcUnbindOxg(1, tstart, tmid) + CalcUnbindOxg(1, tmid, tend)
		assert.Equal(t, total, total2)
	}
}

// test 1 balance will get ONYX_TOTAL_SUPPLY eventually
func TestTotalONG(t *testing.T) {
	assert.Equal(t, CalcUnbindOxg(1, 0, constants.UNBOUND_DEADLINE),
		constants.ONYX_TOTAL_SUPPLY)

	assert.Equal(t, CalcUnbindOxg(1, 0, TIME_INTERVAL*18),
		constants.ONYX_TOTAL_SUPPLY)

	assert.Equal(t, CalcUnbindOxg(1, 0, TIME_INTERVAL*108),
		constants.ONYX_TOTAL_SUPPLY)

	assert.Equal(t, CalcUnbindOxg(1, 0, ^uint32(0)),
		constants.ONYX_TOTAL_SUPPLY)
}
