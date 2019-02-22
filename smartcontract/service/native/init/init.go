/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package init

import (
	"bytes"
	"math/big"

	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/auth"
	params "github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/global_params"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/governance"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/ong"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/ont"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/ontid"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/neovm"
	vm "github.com/OnyxPay/OnyxChain-legacy/vm/neovm"
)

var (
	COMMIT_DPOS_BYTES = InitBytes(utils.GovernanceContractAddress, governance.COMMIT_DPOS)
)

func init() {
	ong.InitOng()
	ont.InitOnt()
	params.InitGlobalParams()
	ontid.Init()
	auth.Init()
	governance.InitGovernance()
}

func InitBytes(addr common.Address, method string) []byte {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	return builder.ToArray()
}
