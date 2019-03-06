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

func TestFormatOnx(t *testing.T) {
	assert.Equal(t, "0", FormatOnx(0))
	assert.Equal(t, "1", FormatOnx(1))
	assert.Equal(t, "100", FormatOnx(100))
	assert.Equal(t, "1000000000", FormatOnx(1000000000))
}

func TestParseOnx(t *testing.T) {
	assert.Equal(t, uint64(0), ParseOnx("0"))
	assert.Equal(t, uint64(1), ParseOnx("1"))
	assert.Equal(t, uint64(1000), ParseOnx("1000"))
	assert.Equal(t, uint64(1000000000), ParseOnx("1000000000"))
	assert.Equal(t, uint64(1000000), ParseOnx("1000000.123"))
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
