package test

import (
	tronWallet "github.com/nanchang0000/tron-wallet"
	"github.com/nanchang0000/tron-wallet/enums"
)

var node = enums.SHASTA_NODE
var validPrivateKey = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
var invalidPrivateKey = "invalid"
var validOwnerAddress = "TSw5FSuWhTAcaJmBUVFY9fUY4ihwx588b6"
var invalidOwnerAddress = "T2w5FSuWhxcaJmBUVFY93UY4ihwx55668b6"
var validToAddress = "TJnsY5bGiwuPCQFismQDwyVTPAn7M88888"
var invalidToAddress = "TJnsY5bGiwuPCQQDwyVTPAnM88888"
var trxAmount int64 = 10000
var trc20Amount int64 = 10000
var nileNode = enums.NILE_NODE
var childAddress = "TNenhn2ennQvxcRuMZ6VKKPacMDDujPjNK"
var childPrivateKey = "cc492f384b5929feed9995abb103a0fc914304ac319dbd3efbb813abcc3dbd9f"
var permissionAddress = "TVv3U8MDTwDzH2RAF47dJNbhMxqzWnFRji"
var permissionPrivateKey = "b4b50da5a86ef17148b6e14ca9a86e71ab09e7724a7445f766c5819ec550ab8f"

func wallet() *tronWallet.TronWallet {
	w, _ := tronWallet.CreateTronWallet(node, validPrivateKey)
	return w
}

func token() *tronWallet.Token {
	return &tronWallet.Token{
		ContractAddress: enums.SHASTA_Tether_USDT,
	}
}

func crawler() *tronWallet.Crawler {
	return &tronWallet.Crawler{
		Node:      node,
		Addresses: []string{validOwnerAddress},
	}
}
