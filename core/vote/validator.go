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

package vote

import (
	"github.com/OnyxPay/OnyxChain-crypto/keypair"
	"github.com/OnyxPay/OnyxChain-legacy/core/genesis"
	"github.com/OnyxPay/OnyxChain-legacy/core/types"
)

func GetValidators(txs []*types.Transaction) ([]keypair.PublicKey, error) {
	// TODO implement vote
	return genesis.GenesisBookkeepers, nil
}