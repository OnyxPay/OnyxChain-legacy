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
package types

import (
	"bytes"
	"testing"

	"github.com/OnyxPay/OnyxChain-legacy/p2pserver/common"
	"github.com/stretchr/testify/assert"
)

func TestMsgHdrSerializationDeserialization(t *testing.T) {
	hdr := newMessageHeader("hdrtest", 0, common.Checksum(nil))

	buf := bytes.NewBuffer(nil)
	err := writeMessageHeader(buf, hdr)
	if err != nil {
		return
	}

	dehdr, err := readMessageHeader(buf)
	assert.Nil(t, err)

	assert.Equal(t, hdr, dehdr)

}
