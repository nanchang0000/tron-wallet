package grpcClient

import (
	"bytes"
	"fmt"
	"github.com/nanchang0000/tron-wallet/grpcClient/proto/api"
	"github.com/nanchang0000/tron-wallet/grpcClient/proto/core"
	"github.com/nanchang0000/tron-wallet/util"
	"google.golang.org/protobuf/proto"
	"math/big"
)

func (g *GrpcClient) GetAccount(addr string) (*core.Account, error) {
	account := new(core.Account)
	var err error

	account.Address, err = util.DecodeCheck(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := g.getContext()
	defer cancel()

	acc, err := g.Client.GetAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(acc.Address, account.Address) {
		return nil, fmt.Errorf("account not found")
	}
	return acc, nil
}

func (g *GrpcClient) GetAccountResource(addr string) (*api.AccountResourceMessage, error) {
	account := new(core.Account)
	var err error

	account.Address, err = util.DecodeCheck(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := g.getContext()
	defer cancel()

	acc, err := g.Client.GetAccountResource(ctx, account)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func makePermission(name string, pType core.Permission_PermissionType, id int32,
	threshold int64, operations map[string]bool, keys map[string]int64) (*core.Permission, error) {

	pKey := make([]*core.Key, 0)

	if len(keys) > 5 {
		return nil, fmt.Errorf("cant have more than 5 keys")
	}
	totalWeight := int64(0)
	for k, w := range keys {
		totalWeight += w
		addr, err := util.Base58ToAddress(k)
		if err != nil {
			return nil, fmt.Errorf("invalid address: %s", k)
		}
		pKey = append(pKey, &core.Key{
			Address: addr,
			Weight:  w,
		})
	}
	var bigOP *big.Int
	if operations != nil && len(operations) > 0 {
		bigOP = big.NewInt(0)
		for k, o := range operations {
			if o {
				// find k in contracts
				value, b := core.Transaction_Contract_ContractType_value[k]
				if !b {
					return nil, fmt.Errorf("permission not found: %s", k)
				}
				bigOP.SetBit(bigOP, int(value), 1)
			}
		}
	} else {
		bigOP = nil
	}

	if threshold > totalWeight {
		return nil, fmt.Errorf("invalid key/threshold size (%d/%d)", threshold, totalWeight)
	}
	var bOP []byte
	if bigOP != nil {
		bOP = make([]byte, 32)
		l := len(bigOP.Bytes()) - 1
		for i, b := range bigOP.Bytes() {
			bOP[l-i] = b
		}
	}

	return &core.Permission{
		Type:           pType,
		Id:             id,
		PermissionName: name,
		Threshold:      threshold,
		Operations:     bOP,
		Keys:           pKey,
	}, nil
}

// UpdateAccountPermission change account permission
func (g *GrpcClient) UpdateAccountPermission(from string, owner, witness map[string]interface{}, actives []map[string]interface{}) (*api.TransactionExtention, error) {

	if len(actives) > 8 {
		return nil, fmt.Errorf("cant have more than 8 active operations")
	}

	if owner == nil {
		return nil, fmt.Errorf("owner is manadory")
	}
	ownerPermission, err := makePermission(
		"owner",
		core.Permission_Owner,
		0,
		owner["threshold"].(int64),
		nil,
		owner["keys"].(map[string]int64),
	)
	if err != nil {
		return nil, err
	}
	contract := &core.AccountPermissionUpdateContract{
		Owner: ownerPermission,
	}

	if contract.OwnerAddress, err = util.DecodeCheck(from); err != nil {
		return nil, err
	}

	if actives != nil {
		activesPermission := make([]*core.Permission, 0)
		for i, active := range actives {
			activeP, err := makePermission(
				active["name"].(string),
				core.Permission_Active,
				int32(2+i),
				active["threshold"].(int64),
				active["operations"].(map[string]bool),
				active["keys"].(map[string]int64),
			)
			if err != nil {
				return nil, err
			}
			activesPermission = append(activesPermission, activeP)
		}
		contract.Actives = activesPermission
	}

	if witness != nil {
		witnessPermission, err := makePermission(
			"witness",
			core.Permission_Witness,
			1,
			witness["threshold"].(int64),
			nil,
			witness["keys"].(map[string]int64),
		)
		if err != nil {
			return nil, err
		}
		contract.Witness = witnessPermission
	}

	ctx, cancel := g.getContext()
	defer cancel()

	tx, err := g.Client.AccountPermissionUpdate(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}
