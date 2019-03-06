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

package db

import (
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/core/types"
)

type BestBlock struct {
	Height uint32
	Hash   common.Uint256
}

type BestStateProvider interface {
	GetBestBlock() (BestBlock, error)
	GetBestHeader() (*types.Header, error)
}
