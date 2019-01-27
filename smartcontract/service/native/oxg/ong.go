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

package oxg

import (
	"math/big"

	"fmt"
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/common/constants"
	"github.com/OnyxPay/OnyxChain-legacy/errors"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/onyx"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/OnyxPay/OnyxChain-legacy/vm/neovm/types"
)

func InitOxg() {
	native.Contracts[utils.OxgContractAddress] = RegisterOxgContract
}

func RegisterOxgContract(native *native.NativeService) {
	native.Register(onyx.INIT_NAME, OxgInit)
	native.Register(onyx.TRANSFER_NAME, OxgTransfer)
	native.Register(onyx.APPROVE_NAME, OxgApprove)
	native.Register(onyx.TRANSFERFROM_NAME, OxgTransferFrom)
	native.Register(onyx.NAME_NAME, OxgName)
	native.Register(onyx.SYMBOL_NAME, OxgSymbol)
	native.Register(onyx.DECIMALS_NAME, OxgDecimals)
	native.Register(onyx.TOTALSUPPLY_NAME, OxgTotalSupply)
	native.Register(onyx.BALANCEOF_NAME, OxgBalanceOf)
	native.Register(onyx.ALLOWANCE_NAME, OxgAllowance)
}

func OxgInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, onyx.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init oxg has been completed!")
	}

	item := utils.GenUInt64StorageItem(constants.OXG_TOTAL_SUPPLY)
	native.CacheDB.Put(onyx.GenTotalSupplyKey(contract), item.ToArray())
	native.CacheDB.Put(append(contract[:], utils.OnyxContractAddress[:]...), item.ToArray())
	onyx.AddNotifications(native, contract, &onyx.State{To: utils.OnyxContractAddress, Value: constants.OXG_TOTAL_SUPPLY})
	return utils.BYTE_TRUE, nil
}

func OxgTransfer(native *native.NativeService) ([]byte, error) {
	var transfers onyx.Transfers
	source := common.NewZeroCopySource(native.Input)
	if err := transfers.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OxgTransfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.OXG_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer oxg amount:%d over totalSupply:%d", v.Value, constants.OXG_TOTAL_SUPPLY)
		}
		if _, _, err := onyx.Transfer(native, contract, &v); err != nil {
			return utils.BYTE_FALSE, err
		}
		onyx.AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func OxgApprove(native *native.NativeService) ([]byte, error) {
	var state onyx.State
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OxgApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.OXG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve oxg amount:%d over totalSupply:%d", state.Value, constants.OXG_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CacheDB.Put(onyx.GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value).ToArray())
	return utils.BYTE_TRUE, nil
}

func OxgTransferFrom(native *native.NativeService) ([]byte, error) {
	var state onyx.TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnyxTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.OXG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve oxg amount:%d over totalSupply:%d", state.Value, constants.OXG_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if _, _, err := onyx.TransferedFrom(native, contract, &state); err != nil {
		return utils.BYTE_FALSE, err
	}
	onyx.AddNotifications(native, contract, &onyx.State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func OxgName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.OXG_NAME), nil
}

func OxgDecimals(native *native.NativeService) ([]byte, error) {
	return big.NewInt(int64(constants.OXG_DECIMALS)).Bytes(), nil
}

func OxgSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.OXG_SYMBOL), nil
}

func OxgTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, onyx.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnyxTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func OxgBalanceOf(native *native.NativeService) ([]byte, error) {
	return onyx.GetBalanceValue(native, onyx.TRANSFER_FLAG)
}

func OxgAllowance(native *native.NativeService) ([]byte, error) {
	return onyx.GetBalanceValue(native, onyx.APPROVE_FLAG)
}
