package tronWallet

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"

	"github.com/nanchang0000/tron-wallet/enums"
	"github.com/nanchang0000/tron-wallet/grpcClient"
)

func (t *TronWallet) CreateAndBroadcastMultiTransaction(node enums.Node, fromAddressBase58 string, toAddressBase58 string, amountInSun int64, privateKeys []*ecdsa.PrivateKey, _ ecdsa.PrivateKey) (string, error) {

	transaction, err := createTransactionInput(node, fromAddressBase58, toAddressBase58, amountInSun)
	if err != nil {
		return "", err
	}

	for _, privateKey := range privateKeys {
		transaction, err = signTransaction(transaction, privateKey)
		if err != nil {
			return "", err
		}
	}

	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return "", err
	}

	res, err := c.Broadcast(transaction.Transaction)
	if err != nil {
		return "", err
	}

	if !res.Result {
		return "", errors.New(res.Code.String())
	}

	return string(transaction.GetTxid()), nil
}

func (t *TronWallet) MultiTransferTrc20(node enums.Node, permissionWallet *TronWallet, targetAddress string, amount *big.Int, token *Token) (string, error) {
	childPrivateKey, _ := t.PrivateKeyRCDSA()
	permissionPrivateKey, _ := permissionWallet.PrivateKeyRCDSA()
	transaction, err := createTrc20TransactionInput(node, t.AddressBase58, token, targetAddress, amount)
	if err != nil {
		return "", err
	}
	transaction, err = signTransaction(transaction, childPrivateKey)
	if err != nil {
		return "", err
	}
	transaction, err = signTransaction(transaction, permissionPrivateKey)
	if err != nil {
		return "", err
	}
	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return "", err
	}
	res, err := c.Broadcast(transaction.Transaction)
	if err != nil {
		return "", err
	}
	if !res.Result {
		return "", errors.New(res.Code.String())
	}
	return hexutil.Encode(transaction.GetTxid())[2:], nil
}
