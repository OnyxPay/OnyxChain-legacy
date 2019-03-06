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

package errors

import (
	"errors"
)

const callStackDepth = 10

type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetRoot() error
}

func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}

	onxerr, ok := err.(onxError)
	if !ok {
		onxerr.root = err
		onxerr.errmsg = err.Error()
		onxerr.callstack = getCallStack(0, callStackDepth)
		onxerr.code = errcode

	}
	if errmsg != "" {
		onxerr.errmsg = errmsg + ": " + onxerr.errmsg
	}

	return onxerr
}

func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetRoot()
	}
	return err
}
