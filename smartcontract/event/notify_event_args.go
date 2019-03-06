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

package event

import (
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/vm/neovm/types"
)

const (
	CONTRACT_STATE_FAIL    byte = 0
	CONTRACT_STATE_SUCCESS byte = 1
)

// NotifyEventArgs describe smart contract event notify arguments struct
type NotifyEventArgs struct {
	ContractAddress common.Address
	States          types.StackItems
}

// NotifyEventInfo describe smart contract event notify info struct
type NotifyEventInfo struct {
	ContractAddress common.Address
	States          interface{}
}

type ExecuteNotify struct {
	TxHash      common.Uint256
	State       byte
	GasConsumed uint64
	Notify      []*NotifyEventInfo
}
