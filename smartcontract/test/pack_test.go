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

package test

import (
	"github.com/OnyxPay/OnyxChain-legacy/core/types"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract"
	"github.com/OnyxPay/OnyxChain-legacy/vm/neovm"
	"os"
	"testing"
)

func TestPackCrash(t *testing.T) {
	// define a leaf
	byteCode := []byte{byte(neovm.PUSH0)}

	// build tree using array packing
	for i := 0; i < 10000; i++ {
		byteCode = append(byteCode, []byte{
			byte(neovm.DUP),
			byte(neovm.PUSH2),
			byte(neovm.PACK),
		}...)
	}
	// compare trees
	byteCode = append(byteCode, []byte{
		byte(neovm.DUP),
		byte(neovm.EQUAL),
	}...)
	// setup VM
	dbFile := "test"
	os.RemoveAll(dbFile)
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	panic(err)
	//}
	//store := statestore.NewMemDatabase()
	//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
	config := &smartcontract.Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	//cache := storage.NewCloneCache(testBatch)
	sc := smartcontract.SmartContract{
		Config:  config,
		Gas:     200,
		CacheDB: nil,
	}
	engine, err := sc.NewExecuteEngine(byteCode)
	if err != nil {
		panic(err)
		// cause the VM to hang forever
		_, err = engine.Invoke()
		if err != nil {
		}
		panic(err)
	}
}
