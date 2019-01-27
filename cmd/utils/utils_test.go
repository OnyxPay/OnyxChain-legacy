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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatOxg(t *testing.T) {
	assert.Equal(t, "1", FormatOxg(1000000000))
	assert.Equal(t, "1.1", FormatOxg(1100000000))
	assert.Equal(t, "1.123456789", FormatOxg(1123456789))
	assert.Equal(t, "1000000000.123456789", FormatOxg(1000000000123456789))
	assert.Equal(t, "1000000000.000001", FormatOxg(1000000000000001000))
	assert.Equal(t, "1000000000.000000001", FormatOxg(1000000000000000001))
}

func TestParseOxg(t *testing.T) {
	assert.Equal(t, uint64(1000000000), ParseOxg("1"))
	assert.Equal(t, uint64(1000000000000000000), ParseOxg("1000000000"))
	assert.Equal(t, uint64(1000000000123456789), ParseOxg("1000000000.123456789"))
	assert.Equal(t, uint64(1000000000000000100), ParseOxg("1000000000.0000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseOxg("1000000000.000000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseOxg("1000000000.000000001123"))
}

func TestFormatOnyx(t *testing.T) {
	assert.Equal(t, "0", FormatOnyx(0))
	assert.Equal(t, "1", FormatOnyx(1))
	assert.Equal(t, "100", FormatOnyx(100))
	assert.Equal(t, "1000000000", FormatOnyx(1000000000))
}

func TestParseOnyx(t *testing.T) {
	assert.Equal(t, uint64(0), ParseOnyx("0"))
	assert.Equal(t, uint64(1), ParseOnyx("1"))
	assert.Equal(t, uint64(1000), ParseOnyx("1000"))
	assert.Equal(t, uint64(1000000000), ParseOnyx("1000000000"))
	assert.Equal(t, uint64(1000000), ParseOnyx("1000000.123"))
}

func TestGenExportBlocksFileName(t *testing.T) {
	name := "blocks.dat"
	start := uint32(0)
	end := uint32(100)
	fileName := GenExportBlocksFileName(name, start, end)
	assert.Equal(t, "blocks_0_100.dat", fileName)
	name = "blocks"
	fileName = GenExportBlocksFileName(name, start, end)
	assert.Equal(t, "blocks_0_100", fileName)
	name = "blocks."
	fileName = GenExportBlocksFileName(name, start, end)
	assert.Equal(t, "blocks_0_100.", fileName)
	name = "blocks.export.dat"
	fileName = GenExportBlocksFileName(name, start, end)
	assert.Equal(t, "blocks.export_0_100.dat", fileName)
}
