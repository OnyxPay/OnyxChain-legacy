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

package cmd

import (
	"fmt"
	"github.com/OnyxPay/OnyxChain-legacy/account"
	cmdcom "github.com/OnyxPay/OnyxChain-legacy/cmd/common"
	"github.com/OnyxPay/OnyxChain-legacy/cmd/utils"
	"github.com/OnyxPay/OnyxChain-legacy/common/config"
	nutils "github.com/OnyxPay/OnyxChain-legacy/smartcontract/service/native/utils"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

var AssetCommand = cli.Command{
	Name:        "asset",
	Usage:       "Handle assets",
	Description: "Asset management commands can check account balance, ONYX/OXG transfers, extract OXGs, and view unbound OXGs, and so on.",
	Subcommands: []cli.Command{
		{
			Action:      transfer,
			Name:        "transfer",
			Usage:       "Transfer onx or oxg to another account",
			ArgsUsage:   " ",
			Description: "Transfer onx or oxg to another account. If from address does not specified, using default account",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionAssetFlag,
				utils.TransactionFromFlag,
				utils.TransactionToFlag,
				utils.TransactionAmountFlag,
				utils.ForceSendTxFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    approve,
			Name:      "approve",
			ArgsUsage: " ",
			Usage:     "Approve another user can transfer asset",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.ApproveAssetFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.ApproveAmountFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    transferFrom,
			Name:      "transferfrom",
			ArgsUsage: " ",
			Usage:     "Using to transfer asset after approve",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.ApproveAssetFlag,
				utils.TransferFromSenderFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.TransferFromAmountFlag,
				utils.ForceSendTxFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    getBalance,
			Name:      "balance",
			Usage:     "Show balance of onx and oxg of specified account",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action: getAllowance,
			Name:   "allowance",
			Usage:  "Show approve balance of onx or oxg of specified account",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.ApproveAssetFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    unboundOxg,
			Name:      "unboundoxg",
			Usage:     "Show the balance of unbound OXG",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    withdrawOxg,
			Name:      "withdrawoxg",
			Usage:     "Withdraw OXG",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.WalletFileFlag,
			},
		},
	},
}

func transfer(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.TransactionToFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionAmountFlag)) {
		PrintErrorMsg("Missing %s %s or %s argument.", utils.TransactionToFlag.Name, utils.TransactionFromFlag.Name, utils.TransactionAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	asset := ctx.String(utils.GetFlagName(utils.TransactionAssetFlag))
	if asset == "" {
		asset = utils.ASSET_ONX
	}
	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	to := ctx.String(utils.TransactionToFlag.Name)
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var amount uint64
	amountStr := ctx.String(utils.TransactionAmountFlag.Name)
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	force := ctx.Bool(utils.GetFlagName(utils.ForceSendTxFlag))
	if !force {
		balance, err := utils.GetAccountBalance(fromAddr, asset)
		if err != nil {
			return err
		}
		if balance < amount {
			PrintErrorMsg("Account:%s balance not enough.", fromAddr)
			PrintInfoMsg("\nTip:")
			PrintInfoMsg("  If you want to send transaction compulsively, please using %s flag.", utils.GetFlagName(utils.ForceSendTxFlag))
			return nil
		}
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return err
	}
	txHash, err := utils.Transfer(gasPrice, gasLimit, signer, asset, fromAddr, toAddr, amount)
	if err != nil {
		return fmt.Errorf("transfer error:%s", err)
	}
	PrintInfoMsg("Transfer %s", strings.ToUpper(asset))
	PrintInfoMsg("  From:%s", fromAddr)
	PrintInfoMsg("  To:%s", toAddr)
	PrintInfoMsg("  Amount:%s", amountStr)
	PrintInfoMsg("  TxHash:%s", txHash)
	PrintInfoMsg("\nTip:")
	PrintInfoMsg("  Using './onyxchain info status %s' to query transaction status.", txHash)
	return nil
}

func getBalance(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		PrintErrorMsg("Missing account argument.")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	balance, err := utils.GetBalance(accAddr)
	if err != nil {
		return err
	}

	oxg, err := strconv.ParseUint(balance.Oxg, 10, 64)
	if err != nil {
		return err
	}
	PrintInfoMsg("BalanceOf:%s", accAddr)
	PrintInfoMsg("  ONYX:%s", balance.Onx)
	PrintInfoMsg("  OXG:%s", utils.FormatOxg(oxg))
	return nil
}

func getAllowance(ctx *cli.Context) error {
	SetRpcPort(ctx)
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	if from == "" || to == "" {
		PrintErrorMsg("Missing %s or %s argument.", utils.ApproveAssetFromFlag.Name, utils.ApproveAssetToFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	if asset == "" {
		asset = utils.ASSET_ONX
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}
	balanceStr, err := utils.GetAllowance(asset, fromAddr, toAddr)
	if err != nil {
		return err
	}
	switch strings.ToLower(asset) {
	case "onyx":
	case "oxg":
		balance, err := strconv.ParseUint(balanceStr, 10, 64)
		if err != nil {
			return err
		}
		balanceStr = utils.FormatOxg(balance)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}
	PrintInfoMsg("Allowance:%s", asset)
	PrintInfoMsg("  From:%s", fromAddr)
	PrintInfoMsg("  To:%s", toAddr)
	PrintInfoMsg("  Balance:%s", balanceStr)
	return nil
}

func approve(ctx *cli.Context) error {
	SetRpcPort(ctx)
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.ApproveAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		PrintErrorMsg("Missing %s %s %s or %s argument.", utils.ApproveAssetFlag.Name, utils.ApproveAssetFromFlag.Name, utils.ApproveAssetToFlag.Name, utils.ApproveAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}
	var amount uint64
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)
	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return err
	}

	txHash, err := utils.Approve(gasPrice, gasLimit, signer, asset, fromAddr, toAddr, amount)
	if err != nil {
		return fmt.Errorf("approve error:%s", err)
	}

	PrintInfoMsg("Approve:")
	PrintInfoMsg("  Asset:%s", asset)
	PrintInfoMsg("  From:%s", fromAddr)
	PrintInfoMsg("  To:%s", toAddr)
	PrintInfoMsg("  Amount:%s", amountStr)
	PrintInfoMsg("  TxHash:%s", txHash)
	PrintInfoMsg("\nTip:")
	PrintInfoMsg("  Using './onyxchain info status %s' to query transaction status.", txHash)
	return nil
}

func transferFrom(ctx *cli.Context) error {
	SetRpcPort(ctx)
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.TransferFromAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		PrintErrorMsg("Missing %s %s %s or %s argument.", utils.ApproveAssetFlag.Name, utils.ApproveAssetFromFlag.Name, utils.ApproveAssetToFlag.Name, utils.TransferFromAmountFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var sendAddr string
	sender := ctx.String(utils.GetFlagName(utils.TransferFromSenderFlag))
	if sender == "" {
		sendAddr = toAddr
	} else {
		sendAddr, err = cmdcom.ParseAddress(sender, ctx)
		if err != nil {
			return err
		}
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, sendAddr)
	if err != nil {
		return err
	}

	var amount uint64
	switch strings.ToLower(asset) {
	case "onyx":
		amount = utils.ParseOnx(amountStr)
		amountStr = utils.FormatOnx(amount)
	case "oxg":
		amount = utils.ParseOxg(amountStr)
		amountStr = utils.FormatOxg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	force := ctx.Bool(utils.GetFlagName(utils.ForceSendTxFlag))
	if !force {
		balance, err := utils.GetAccountBalance(fromAddr, asset)
		if err != nil {
			return err
		}
		if balance < amount {
			PrintErrorMsg("Account:%s balance not enough.", fromAddr)
			PrintInfoMsg("\nTip:")
			PrintInfoMsg("  If you want to send transaction compulsively, please using %s flag.", utils.GetFlagName(utils.ForceSendTxFlag))
			return nil
		}
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)
	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	txHash, err := utils.TransferFrom(gasPrice, gasLimit, signer, asset, sendAddr, fromAddr, toAddr, amount)
	if err != nil {
		return err
	}

	PrintInfoMsg("Transfer from:")
	PrintInfoMsg("  Asset:%s", asset)
	PrintInfoMsg("  Sender:%s", sendAddr)
	PrintInfoMsg("  From:%s", fromAddr)
	PrintInfoMsg("  To:%s", toAddr)
	PrintInfoMsg("  Amount:%s", amountStr)
	PrintInfoMsg("  TxHash:%s", txHash)
	PrintInfoMsg("\nTip:")
	PrintInfoMsg("  Using './onyxchain info status %s' to query transaction status.", txHash)
	return nil
}

func unboundOxg(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		PrintErrorMsg("Missing account argument.")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	fromAddr := nutils.OnxContractAddress.ToBase58()
	balanceStr, err := utils.GetAllowance("oxg", fromAddr, accAddr)
	if err != nil {
		return err
	}
	balance, err := strconv.ParseUint(balanceStr, 10, 64)
	if err != nil {
		return err
	}
	balanceStr = utils.FormatOxg(balance)
	PrintInfoMsg("Unbound OXG:")
	PrintInfoMsg("  Account:%s", accAddr)
	PrintInfoMsg("  OXG:%s", balanceStr)
	return nil
}

func withdrawOxg(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		PrintErrorMsg("Missing account argument.")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	fromAddr := nutils.OnxContractAddress.ToBase58()
	balance, err := utils.GetAllowance("oxg", fromAddr, accAddr)
	if err != nil {
		return err
	}

	amount, err := strconv.ParseUint(balance, 10, 64)
	if err != nil {
		return err
	}
	if amount <= 0 {
		return fmt.Errorf("haven't unbound oxg\n")
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, accAddr)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)
	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	txHash, err := utils.TransferFrom(gasPrice, gasLimit, signer, "oxg", accAddr, fromAddr, accAddr, amount)
	if err != nil {
		return err
	}

	PrintInfoMsg("Withdraw OXG:")
	PrintInfoMsg("  Account:%s", accAddr)
	PrintInfoMsg("  Amount:%s", utils.FormatOxg(amount))
	PrintInfoMsg("  TxHash:%s", txHash)
	PrintInfoMsg("\nTip:")
	PrintInfoMsg("  Using './onyxchain info status %s' to query transaction status.", txHash)
	return nil
}
