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
package onxid

import (
	"encoding/hex"

	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/event"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native"
)

func newEvent(srvc *native.NativeService, st interface{}) {
	e := event.NotifyEventInfo{}
	e.ContractAddress = srvc.ContextRef.CurrentContext().ContractAddress
	e.States = st
	srvc.Notifications = append(srvc.Notifications, &e)
	return
}

func triggerRegisterEvent(srvc *native.NativeService, id []byte) {
	newEvent(srvc, []string{"Register", string(id)})
}

func triggerPublicEvent(srvc *native.NativeService, op string, id, pub []byte, keyID uint32) {
	st := []interface{}{"PublicKey", op, string(id), keyID, hex.EncodeToString(pub)}
	newEvent(srvc, st)
}

func triggerAttributeEvent(srvc *native.NativeService, op string, id []byte, path [][]byte) {
	var attr interface{}
	if op == "remove" {
		attr = hex.EncodeToString(path[0])
	} else {
		t := make([]string, len(path))
		for i, v := range path {
			t[i] = hex.EncodeToString(v)
		}
		attr = t
	}
	st := []interface{}{"Attribute", op, string(id), attr}
	newEvent(srvc, st)
}

func triggerRecoveryEvent(srvc *native.NativeService, op string, id []byte, addr common.Address) {
	st := []string{"Recovery", op, string(id), addr.ToHexString()}
	newEvent(srvc, st)
}
