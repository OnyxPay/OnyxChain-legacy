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
package wasmvm

import (
	"bytes"
	"github.com/OnyxPay/OnyxChain-legacy/core/types"
	"github.com/OnyxPay/OnyxChain-legacy/errors"
	"github.com/OnyxPay/OnyxChain-legacy/vm/wasmvm/exec"
)

func (this *WasmVmService) attributeGetUsage(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[transactionGetHash] parameter count error")
	}

	attributebytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, nil
	}

	attr := types.TxAttribute{}
	err = attr.Deserialize(bytes.NewBuffer(attributebytes))
	if err != nil {
		return false, nil
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(attr.Usage))
	return true, nil
}
func (this *WasmVmService) attributeGetData(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[transactionGetHash] parameter count error")
	}

	attributebytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, nil
	}

	attr := types.TxAttribute{}
	err = attr.Deserialize(bytes.NewBuffer(attributebytes))
	if err != nil {
		return false, nil
	}

	idx, err := vm.SetPointerMemory(attr.Data)
	if err != nil {
		return false, nil
	}

	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil
}
