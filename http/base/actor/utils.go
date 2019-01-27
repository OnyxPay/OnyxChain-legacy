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

// Package actor privides communication with other actor
package actor

import (
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
)

func updateNativeSCAddr(hash common.Address) common.Address {
	if hash == utils.OnyxContractAddress {
		hash = common.AddressFromVmCode(utils.OnyxContractAddress[:])
	} else if hash == utils.OxgContractAddress {
		hash = common.AddressFromVmCode(utils.OxgContractAddress[:])
	} else if hash == utils.OnyxIDContractAddress {
		hash = common.AddressFromVmCode(utils.OnyxIDContractAddress[:])
	} else if hash == utils.ParamContractAddress {
		hash = common.AddressFromVmCode(utils.ParamContractAddress[:])
	} else if hash == utils.AuthContractAddress {
		hash = common.AddressFromVmCode(utils.AuthContractAddress[:])
	} else if hash == utils.GovernanceContractAddress {
		hash = common.AddressFromVmCode(utils.GovernanceContractAddress[:])
	}
	return hash
}
