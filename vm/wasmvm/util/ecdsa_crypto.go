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

package util

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/OnyxPay/OnyxChain-crypto/keypair"
	s "github.com/OnyxPay/OnyxChain-crypto/signature"
	"github.com/OnyxPay/OnyxChain-legacy/common/log"
	onxErrors "github.com/OnyxPay/OnyxChain-legacy/errors"
	"golang.org/x/crypto/ripemd160"
)

type ECDsaCrypto struct {
}

func (c *ECDsaCrypto) Hash160(message []byte) []byte {
	temp := sha256.Sum256(message)
	md := ripemd160.New()
	io.WriteString(md, string(temp[:]))
	hash := md.Sum(nil)
	return hash
}

func (c *ECDsaCrypto) Hash256(message []byte) []byte {
	temp := sha256.Sum256(message)
	f := sha256.Sum256(temp[:])
	return f[:]
}

func (c *ECDsaCrypto) VerifySignature(message []byte, signature []byte, pubkey []byte) (bool, error) {

	log.Debugf("message: %x", message)
	log.Debugf("signature: %x", signature)
	log.Debugf("pubkey: %x", pubkey)

	pk, err := keypair.DeserializePublicKey(pubkey)
	if err != nil {
		return false, onxErrors.NewDetailErr(errors.New("[ECDsaCrypto], deserializing public key failed."), onxErrors.ErrNoCode, "")
	}

	sig, err := s.Deserialize(signature)
	ok := s.Verify(pk, message, sig)
	if !ok {
		return false, onxErrors.NewDetailErr(errors.New("[ECDsaCrypto], VerifySignature failed."), onxErrors.ErrNoCode, "")
	}

	return true, nil
}
