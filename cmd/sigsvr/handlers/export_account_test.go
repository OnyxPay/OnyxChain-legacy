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
package handlers

import (
	"encoding/json"
	"github.com/OnyxPay/OnyxChain-legacy/account"
	clisvrcom "github.com/OnyxPay/OnyxChain-legacy/cmd/sigsvr/common"
	"os"
	"testing"
)

func TestExportWallet(t *testing.T) {
	exportReq := &ExportAccountReq{}
	data, _ := json.Marshal(exportReq)
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "exportaccount",
		Pwd:    string(pwd),
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	ExportAccount(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("ExportAccount failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
	exportRsp, ok := resp.Result.(*ExportAccountResp)
	if !ok {
		t.Errorf("TestExportWallet resp asset to ExportAccountResp failed")
		return
	}

	wallet, err := account.Open(exportRsp.WalletFile)
	if err != nil {
		t.Errorf("TestExportWallet failed, OpenWallet error:%s", err)
		return
	}
	if wallet.GetAccountNum() != exportRsp.AccountNumber {
		t.Errorf("TestExportWallet failed, account number %d != %d", wallet.GetAccountNum(), exportRsp.AccountNumber)
		return
	}
	os.Remove(exportRsp.WalletFile)
}
