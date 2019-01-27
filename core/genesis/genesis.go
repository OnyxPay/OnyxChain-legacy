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

package genesis

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/OnyxPay/OnyxChain-crypto/keypair"
	"github.com/OnyxPay/OnyxChain-legacy/common"
	"github.com/OnyxPay/OnyxChain-legacy/common/config"
	"github.com/OnyxPay/OnyxChain-legacy/common/constants"
	"github.com/OnyxPay/OnyxChain-legacy/consensus/vbft/config"
	"github.com/OnyxPay/OnyxChain-legacy/core/types"
	"github.com/OnyxPay/OnyxChain-legacy/core/utils"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/global_params"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/governance"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/onyx"
	nutils "github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/neovm"
)

const (
	BlockVersion uint32 = 0
	GenesisNonce uint64 = 2083236893
)

var (
	ONTToken   = newGoverningToken()
	ONGToken   = newUtilityToken()
	ONTTokenID = ONTToken.Hash()
	ONGTokenID = ONGToken.Hash()
)

var GenBlockTime = (config.DEFAULT_GEN_BLOCK_TIME * time.Second)

var INIT_PARAM = map[string]string{
	"gasPrice": "0",
}

var GenesisBookkeepers []keypair.PublicKey

// BuildGenesisBlock returns the genesis block with default consensus bookkeeper list
func BuildGenesisBlock(defaultBookkeeper []keypair.PublicKey, genesisConfig *config.GenesisConfig) (*types.Block, error) {
	//getBookkeeper
	GenesisBookkeepers = defaultBookkeeper
	nextBookkeeper, err := types.AddressFromBookkeepers(defaultBookkeeper)
	if err != nil {
		return nil, fmt.Errorf("[Block],BuildGenesisBlock err with GetBookkeeperAddress: %s", err)
	}
	conf := bytes.NewBuffer(nil)
	if genesisConfig.VBFT != nil {
		genesisConfig.VBFT.Serialize(conf)
	}
	govConfig := newGoverConfigInit(conf.Bytes())
	consensusPayload, err := vconfig.GenesisConsensusPayload(govConfig.Hash(), 0)
	if err != nil {
		return nil, fmt.Errorf("consensus genesus init failed: %s", err)
	}
	//blockdata
	genesisHeader := &types.Header{
		Version:          BlockVersion,
		PrevBlockHash:    common.Uint256{},
		TransactionsRoot: common.Uint256{},
		Timestamp:        constants.GENESIS_BLOCK_TIMESTAMP,
		Height:           uint32(0),
		ConsensusData:    GenesisNonce,
		NextBookkeeper:   nextBookkeeper,
		ConsensusPayload: consensusPayload,

		Bookkeepers: nil,
		SigData:     nil,
	}

	//block
	onyx := newGoverningToken()
	oxg := newUtilityToken()
	param := newParamContract()
	oid := deployOnyxIDContract()
	auth := deployAuthContract()
	config := newConfig()

	genesisBlock := &types.Block{
		Header: genesisHeader,
		Transactions: []*types.Transaction{
			onyx,
			oxg,
			param,
			oid,
			auth,
			config,
			newGoverningInit(),
			newUtilityInit(),
			newParamInit(),
			govConfig,
		},
	}
	genesisBlock.RebuildMerkleRoot()
	return genesisBlock, nil
}

func newGoverningToken() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.OnyxContractAddress[:], "ONYX", "1.0",
		"OnyxChain Team", "contact@onyx.io", "OnyxChain Network ONYX Token", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis governing token transaction error ")
	}
	return tx
}

func newUtilityToken() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.OxgContractAddress[:], "oxg", "1.0",
		"OnyxChain Team", "contact@onyx.io", "OnyxChain Network oxg Token", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis utility token transaction error ")
	}
	return tx
}

func newParamContract() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.ParamContractAddress[:],
		"ParamConfig", "1.0", "OnyxChain Team", "contact@onyx.io",
		"Chain Global Environment Variables Manager ", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis param transaction error ")
	}
	return tx
}

func newConfig() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.GovernanceContractAddress[:], "CONFIG", "1.0",
		"OnyxChain Team", "contact@onyx.io", "OnyxChain Network Consensus Config", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis config transaction error ")
	}
	return tx
}

func deployAuthContract() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.AuthContractAddress[:], "AuthContract", "1.0",
		"OnyxChain Team", "contact@onyx.io", "OnyxChain Network Authorization Contract", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis auth transaction error ")
	}
	return tx
}

func deployOnyxIDContract() *types.Transaction {
	mutable := utils.NewDeployTransaction(nutils.OnyxIDContractAddress[:], "OID", "1.0",
		"OnyxChain Team", "contact@onyx.io", "OnyxChain Network ONYX ID", true)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis onyxid transaction error ")
	}
	return tx
}

func newGoverningInit() *types.Transaction {
	bookkeepers, _ := config.DefConfig.GetBookkeepers()

	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		m := (5*len(bookkeepers) + 6) / 7
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}

	distribute := []struct {
		addr  common.Address
		value uint64
	}{{addr, constants.ONYX_TOTAL_SUPPLY}}

	args := bytes.NewBuffer(nil)
	nutils.WriteVarUint(args, uint64(len(distribute)))
	for _, part := range distribute {
		nutils.WriteAddress(args, part.addr)
		nutils.WriteVarUint(args, part.value)
	}

	mutable := utils.BuildNativeTransaction(nutils.OnyxContractAddress, onyx.INIT_NAME, args.Bytes())
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis governing token transaction error ")
	}
	return tx
}

func newUtilityInit() *types.Transaction {
	mutable := utils.BuildNativeTransaction(nutils.OxgContractAddress, onyx.INIT_NAME, []byte{})
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis governing token transaction error ")
	}

	return tx
}

func newParamInit() *types.Transaction {
	params := new(global_params.Params)
	var s []string
	for k := range INIT_PARAM {
		s = append(s, k)
	}

	for k, v := range neovm.INIT_GAS_TABLE {
		INIT_PARAM[k] = strconv.FormatUint(v, 10)
		s = append(s, k)
	}

	sort.Strings(s)
	for _, v := range s {
		params.SetParam(global_params.Param{Key: v, Value: INIT_PARAM[v]})
	}
	bf := new(bytes.Buffer)
	params.Serialize(bf)

	bookkeepers, _ := config.DefConfig.GetBookkeepers()
	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		m := (5*len(bookkeepers) + 6) / 7
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}
	nutils.WriteAddress(bf, addr)

	mutable := utils.BuildNativeTransaction(nutils.ParamContractAddress, global_params.INIT_NAME, bf.Bytes())
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis governing token transaction error ")
	}
	return tx
}

func newGoverConfigInit(config []byte) *types.Transaction {
	mutable := utils.BuildNativeTransaction(nutils.GovernanceContractAddress, governance.INIT_CONFIG, config)
	tx, err := mutable.IntoImmutable()
	if err != nil {
		panic("constract genesis governing token transaction error ")
	}
	return tx
}
