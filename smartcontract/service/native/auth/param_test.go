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
package auth

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
)

var (
	admin    []byte
	newAdmin []byte
	p1       []byte
	p2       []byte
)
var (
	funcs           = []string{"foo1", "foo2"}
	role            = "role"
	OnxContractAddr = utils.OnxContractAddress
)

func init() {
	admin = make([]byte, 32)
	newAdmin = make([]byte, 32)
	p1 = make([]byte, 20)
	p2 = make([]byte, 20)
	rand.Read(admin)
	rand.Read(newAdmin)
	rand.Read(p1)
	rand.Read(p2)
}
func TestSerialization_Init(t *testing.T) {
	param := &InitContractAdminParam{
		AdminOnxID: admin,
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())

	param2 := new(InitContractAdminParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(param.AdminOnxID, param2.AdminOnxID) != 0 {
		t.Fatalf("failed")
	}
}

func TestSerialization_Transfer(t *testing.T) {
	param := &TransferParam{
		ContractAddr:  OnxContractAddr,
		NewAdminOnxID: newAdmin,
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())

	param2 := new(TransferParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, param, param2)
}

func TestSerialization_AssignFuncs(t *testing.T) {
	param := &FuncsToRoleParam{
		ContractAddr: OnxContractAddr,
		AdminOnxID:   admin,
		Role:         []byte("role"),
		FuncNames:    funcs,
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())

	param2 := new(FuncsToRoleParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, param, param2)
}

func TestSerialization_AssignOnxIDs(t *testing.T) {
	param := &OnxIDsToRoleParam{
		ContractAddr: OnxContractAddr,
		AdminOnxID:   admin,
		Role:         []byte(role),
		Persons:      [][]byte{[]byte{0x03, 0x04, 0x05, 0x06}, []byte{0x07, 0x08, 0x09, 0x0a}},
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())
	param2 := new(OnxIDsToRoleParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, param, param2)
}

func TestSerialization_Delegate(t *testing.T) {
	param := &DelegateParam{
		ContractAddr: OnxContractAddr,
		From:         p1,
		To:           p2,
		Role:         []byte(role),
		Period:       60 * 60 * 24,
		Level:        3,
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())
	param2 := new(DelegateParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, param, param2)
}

func TestSerialization_Withdraw(t *testing.T) {
	param := &WithdrawParam{
		ContractAddr: OnxContractAddr,
		Initiator:    p1,
		Delegate:     p2,
		Role:         []byte(role),
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())
	param2 := new(WithdrawParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, param, param2)
}

func TestSerialization_VerifyToken(t *testing.T) {
	param := &VerifyTokenParam{
		ContractAddr: OnxContractAddr,
		Caller:       p1,
		Fn:           "foo1",
	}
	bf := new(bytes.Buffer)
	if err := param.Serialize(bf); err != nil {
		t.Fatal(err)
	}
	rd := bytes.NewReader(bf.Bytes())
	param2 := new(VerifyTokenParam)
	if err := param2.Deserialize(rd); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, param, param2)
}
