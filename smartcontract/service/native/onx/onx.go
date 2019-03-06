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

package onx

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/common/constants"
	"github.com/OnyxPay/OnyxChain-legacy/common/log"
	"github.com/OnyxPay/OnyxChain-legacy/errors"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/OnyxPay/OnyxChain-legacy/vm/neovm/types"
)

const (
	TRANSFER_FLAG byte = 1
	APPROVE_FLAG  byte = 2
)

func InitOnx() {
	native.Contracts[utils.OnxContractAddress] = RegisterOnxContract
}

func RegisterOnxContract(native *native.NativeService) {
	native.Register(INIT_NAME, OnxInit)
	native.Register(TRANSFER_NAME, OnxTransfer)
	native.Register(APPROVE_NAME, OnxApprove)
	native.Register(TRANSFERFROM_NAME, OnxTransferFrom)
	native.Register(NAME_NAME, OnxName)
	native.Register(SYMBOL_NAME, OnxSymbol)
	native.Register(DECIMALS_NAME, OnxDecimals)
	native.Register(TOTALSUPPLY_NAME, OnxTotalSupply)
	native.Register(BALANCEOF_NAME, OnxBalanceOf)
	native.Register(ALLOWANCE_NAME, OnxAllowance)
}

func OnxInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init onx has been completed!")
	}

	distribute := make(map[common.Address]uint64)
	source := common.NewZeroCopySource(native.Input)
	buf, _, irregular, eof := source.NextVarBytes()
	if eof {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadVarBytes, contract params deserialize error!")
	}
	if irregular {
		return utils.BYTE_FALSE, common.ErrIrregularData
	}
	input := common.NewZeroCopySource(buf)
	num, err := utils.DecodeVarUint(input)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("read number error:%v", err)
	}
	sum := uint64(0)
	overflow := false
	for i := uint64(0); i < num; i++ {
		addr, err := utils.DecodeAddress(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read address error:%v", err)
		}
		value, err := utils.DecodeVarUint(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read value error:%v", err)
		}
		sum, overflow = common.SafeAdd(sum, value)
		if overflow {
			return utils.BYTE_FALSE, errors.NewErr("wrong config. overflow detected")
		}
		distribute[addr] += value
	}
	if sum != constants.ONX_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("wrong config. total supply %d != %d", sum, constants.ONX_TOTAL_SUPPLY)
	}

	for addr, val := range distribute {
		balanceKey := GenBalanceKey(contract, addr)
		item := utils.GenUInt64StorageItem(val)
		native.CacheDB.Put(balanceKey, item.ToArray())
		AddNotifications(native, contract, &State{To: addr, Value: val})
	}
	native.CacheDB.Put(GenTotalSupplyKey(contract), utils.GenUInt64StorageItem(constants.ONX_TOTAL_SUPPLY).ToArray())

	return utils.BYTE_TRUE, nil
}

func OnxTransfer(native *native.NativeService) ([]byte, error) {
	var transfers Transfers
	source := common.NewZeroCopySource(native.Input)
	if err := transfers.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[Transfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.ONX_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer onx amount:%d over totalSupply:%d", v.Value, constants.ONX_TOTAL_SUPPLY)
		}
		fromBalance, toBalance, err := Transfer(native, contract, &v)
		if err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantOxg(native, contract, v.From, fromBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantOxg(native, contract, v.To, toBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func OnxTransferFrom(native *native.NativeService) ([]byte, error) {
	var state TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnxTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ONX_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("transferFrom onx amount:%d over totalSupply:%d", state.Value, constants.ONX_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	fromBalance, toBalance, err := TransferedFrom(native, contract, &state)
	if err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantOxg(native, contract, state.From, fromBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantOxg(native, contract, state.To, toBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	AddNotifications(native, contract, &State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func OnxApprove(native *native.NativeService) ([]byte, error) {
	var state State
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OxgApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ONX_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve onx amount:%d over totalSupply:%d", state.Value, constants.ONX_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CacheDB.Put(GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value).ToArray())
	return utils.BYTE_TRUE, nil
}

func OnxName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ONX_NAME), nil
}

func OnxDecimals(native *native.NativeService) ([]byte, error) {
	return types.BigIntToBytes(big.NewInt(int64(constants.ONX_DECIMALS))), nil
}

func OnxSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ONX_SYMBOL), nil
}

func OnxTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OnxTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func OnxBalanceOf(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, TRANSFER_FLAG)
}

func OnxAllowance(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, APPROVE_FLAG)
}

func GetBalanceValue(native *native.NativeService, flag byte) ([]byte, error) {
	source := common.NewZeroCopySource(native.Input)
	from, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	var key []byte
	if flag == APPROVE_FLAG {
		to, err := utils.DecodeAddress(source)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
		}
		key = GenApproveKey(contract, from, to)
	} else if flag == TRANSFER_FLAG {
		key = GenBalanceKey(contract, from)
	}
	amount, err := utils.GetStorageUInt64(native, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] address parse error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func grantOxg(native *native.NativeService, contract, address common.Address, balance uint64) error {
	startOffset, err := getUnboundOffset(native, contract, address)
	if err != nil {
		return err
	}
	if native.Time <= constants.GENESIS_BLOCK_TIMESTAMP {
		return nil
	}
	endOffset := native.Time - constants.GENESIS_BLOCK_TIMESTAMP
	if endOffset < startOffset {
		errstr := fmt.Sprintf("grant Oxg error: wrong timestamp endOffset: %d < startOffset: %d", endOffset, startOffset)
		log.Error(errstr)
		return errors.NewErr(errstr)
	} else if endOffset == startOffset {
		return nil
	}

	if balance != 0 {
		value := utils.CalcUnbindOxg(balance, startOffset, endOffset)

		args, err := getApproveArgs(native, contract, utils.OxgContractAddress, address, value)
		if err != nil {
			return err
		}

		if _, err := native.NativeCall(utils.OxgContractAddress, "approve", args); err != nil {
			return err
		}
	}

	native.CacheDB.Put(genAddressUnboundOffsetKey(contract, address), utils.GenUInt32StorageItem(endOffset).ToArray())
	return nil
}

func getApproveArgs(native *native.NativeService, contract, oxgContract, address common.Address, value uint64) ([]byte, error) {
	bf := new(bytes.Buffer)
	approve := State{
		From:  contract,
		To:    address,
		Value: value,
	}

	stateValue, err := utils.GetStorageUInt64(native, GenApproveKey(oxgContract, approve.From, approve.To))
	if err != nil {
		return nil, err
	}

	approve.Value += stateValue

	if err := approve.Serialize(bf); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}
