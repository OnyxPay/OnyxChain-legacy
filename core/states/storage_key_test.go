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
package states

import (
	"testing"

	"bytes"
	"crypto/rand"

	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/stretchr/testify/assert"
)

func TestStorageKey_Deserialize_Serialize(t *testing.T) {
	var addr common.Address
	rand.Read(addr[:])

	storage := StorageKey{
		ContractAddress: addr,
		Key:             []byte{1, 2, 3},
	}

	buf := bytes.NewBuffer(nil)
	storage.Serialize(buf)
	bs := buf.Bytes()

	var storage2 StorageKey
	storage2.Deserialize(buf)
	assert.Equal(t, storage, storage2)

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err := storage2.Deserialize(buf)
	assert.NotNil(t, err)
}
