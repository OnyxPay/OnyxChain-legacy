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

package neovm

import (
	"math/big"
	"testing"

	vtypes "github.com/OnyxPay/OnyxChain-legacy/vm/neovm/types"
)

func TestOpInvert(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	e.EvaluationStack = stack

	opInvert(&e)
	i := big.NewInt(123456789)

	v, err := PeekBigInteger(&e)
	if err != nil {
		t.Fatal("NeoVM OpInvert test failed.")
	}
	if v.Cmp(i.Not(i)) != 0 {
		t.Fatal("NeoVM OpInvert test failed.")
	}
}

func TestOpEqual(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	e.EvaluationStack = stack

	opEqual(&e)
	v, err := PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpEqual test failed.")
	}
	if !v {
		t.Fatal("NeoVM OpEqual test failed.")
	}
}
