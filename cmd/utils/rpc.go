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

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/OnyxPay/OnyxChain-legacy/common/config"
	rpcerr "github.com/OnyxPay/OnyxChain-legacy/http/base/error"
	"io/ioutil"
	"net/http"
	"strings"
)

//JsonRpc version
const JSON_RPC_VERSION = "2.0"

const (
	ERROR_INVALID_PARAMS   = rpcerr.INVALID_PARAMS
	ERROR_ONYXCHAIN_COMMON  = 10000
	ERROR_ONYXCHAIN_SUCCESS = 0
)

type OnyxChainError struct {
	ErrorCode int64
	Error     error
}

func NewOnyxChainError(err error, errCode ...int64) *OnyxChainError {
	onxErr := &OnyxChainError{Error: err}
	if len(errCode) > 0 {
		onxErr.ErrorCode = errCode[0]
	} else {
		onxErr.ErrorCode = ERROR_ONYXCHAIN_COMMON
	}
	if err == nil {
		onxErr.ErrorCode = ERROR_ONYXCHAIN_SUCCESS
	}
	return onxErr
}

//JsonRpcRequest object in rpc
type JsonRpcRequest struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//JsonRpcResponse object response for JsonRpcRequest
type JsonRpcResponse struct {
	Error  int64           `json:"error"`
	Desc   string          `json:"desc"`
	Result json.RawMessage `json:"result"`
}

func sendRpcRequest(method string, params []interface{}) ([]byte, *OnyxChainError) {
	rpcReq := &JsonRpcRequest{
		Version: JSON_RPC_VERSION,
		Id:      "cli",
		Method:  method,
		Params:  params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, NewOnyxChainError(fmt.Errorf("JsonRpcRequest json.Marshal error:%s", err))
	}

	addr := fmt.Sprintf("http://localhost:%d", config.DefConfig.Rpc.HttpJsonPort)
	resp, err := http.Post(addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, NewOnyxChainError(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewOnyxChainError(fmt.Errorf("read rpc response body error:%s", err))
	}
	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, NewOnyxChainError(fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err))
	}
	if rpcRsp.Error != 0 {
		return nil, NewOnyxChainError(fmt.Errorf("%s", strings.ToLower(rpcRsp.Desc)), rpcRsp.Error)
	}
	return rpcRsp.Result, nil
}
