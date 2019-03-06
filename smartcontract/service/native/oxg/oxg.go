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

package oxg

import (
	"math/big"

	"fmt"
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/common/constants"
	"github.com/OnyxPay/OnyxChain-legacy/errors"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/onx"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/OnyxPay/OnyxChain-legacy/vm/neovm/types"
)

func InitOxg() {
	native.Contracts[utils.OxgContractAddress] = RegisterOxgContract
}

func RegisterOxgContract(native *native.NativeService) {
	native.Register(onx.INIT_NAME, OxgInit)
	native.Register(onx.TRANSFER_NAME, OxgTransfer)
	native.Register(onx.APPROVE_NAME, OxgApprove)
	native.Register(onx.TRANSFERFROM_NAME, OxgTransferFrom)
	native.Register(onx.NAME_NAME, OxgName)
	native.Register(onx.SYMBOL_NAME, OxgSymbol)
	native.Register(onx.DECIMALS_NAME, OxgDecimals)
	native.Register(onx.TOTALSUPPLY_NAME, OxgTotalSupply)
	native.Register(onx.BALANCEOF_NAME, OxgBalanceOf)
	native.Register(onx.ALLOWANCE_NAME, OxgAllowance)
}

func OxgInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, onx.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init oxg has been completed!")
	}

	item := utils.GenUInt64StorageItem(constants.OXG_TOTAL_SUPPLY)
	native.CacheDB.Put(onx.GenTotalSupplyKey(contract), item.ToArray())
	native.CacheDB.Put(append(contract[:], utils.OnxContractAddress[:]...), item.ToArray())
	onx.AddNotifications(native, contract, &onx.State{To: utils.OnxContractAddress, Value: constants.OXG_TOTAL_SUPPLY})
	return utils.BYTE_TRUE, nil
}

func OxgTransfer(native *native.NativeService) ([]byte, error) {
	var transfers onx.Transfers
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
		if _, _, err := onx.Transfer(native, contract, &v); err != nil {
			return utils.BYTE_FALSE, err
		}
		onx.AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func OxgApprove(native *native.NativeService) ([]byte, error) {
	var state onx.State
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
	native.CacheDB.Put(onx.GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value).ToArray())
	return utils.BYTE_TRUE, nil
}

func OxgTransferFrom(native *native.NativeService) ([]byte, error) {
	var state onx.TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnxTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.OXG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve oxg amount:%d over totalSupply:%d", state.Value, constants.OXG_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if _, _, err := onx.TransferedFrom(native, contract, &state); err != nil {
		return utils.BYTE_FALSE, err
	}
	onx.AddNotifications(native, contract, &onx.State{From: state.From, To: state.To, Value: state.Value})
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
	amount, err := utils.GetStorageUInt64(native, onx.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnxTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func OxgBalanceOf(native *native.NativeService) ([]byte, error) {
	return onx.GetBalanceValue(native, onx.TRANSFER_FLAG)
}

func OxgAllowance(native *native.NativeService) ([]byte, error) {
	return onx.GetBalanceValue(native, onx.APPROVE_FLAG)
}
