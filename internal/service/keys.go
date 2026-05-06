package service

import (
	"github.com/gin-gonic/gin"
	"lbe_crypto_signer/pkg/common/apistruct"
	"lbe_crypto_signer/pkg/tools/log"
)

// AddMnemonic signer machine add Mnemonic
func (c *LBESignerT) AddMnemonic(ctx *gin.Context, req *apistruct.KeysAddMnemonicReq) (*apistruct.KeysAddMnemonicResp, error) {
	var (
		err  error
		resp = &apistruct.KeysAddMnemonicResp{}
	)
	mastery, err := c.WalletCli.SeedToMasterKey(ctx, req.Mnemonic, "")
	if err != nil {
		log.ZError(ctx, "AddMnemonic SeedToMasterKey", err)
		return nil, err
	}
	adds, err := c.WalletCli.AccountsFromDerive(ctx, mastery, req.Mnemonic, req.Name)
	if err != nil {
		log.ZError(ctx, "AddMnemonic AccountsFromDerive", err)
		return nil, err
	}
	var accounts []apistruct.AccountBasic
	for _, v := range adds {
		t := apistruct.AccountBasic{
			Addr:  v.Addr,
			Chain: v.Chain,
		}
		accounts = append(accounts, t)
	}
	resp.Accounts = accounts
	return resp, err
}

// KeyList signer machine list key
func (c *LBESignerT) KeyList(ctx *gin.Context, req *apistruct.KeysListReq) (*apistruct.KeysListResp, error) {
	var (
		resp = &apistruct.KeysListResp{}
	)
	keys, err := c.WalletCli.KeyList()
	if err != nil {
		log.ZError(ctx, "KeyList WalletCli.KeyList", err)
		return nil, err
	}
	for _, v := range keys {
		t := apistruct.KeyMeta{}
		t.KeyFLag = v.Key
		t.Name = v.AccountsComp.Name
		for _, a := range v.AccountsComp.Accounts {
			mid := apistruct.AccountBasic{}
			mid.Addr = a.Addr
			mid.Chain = a.Chain
			t.Accounts = append(t.Accounts, mid)
		}
		resp.KeySlice = append(resp.KeySlice, t)
	}
	return resp, err
}

// AddAccounts expand address by index
func (c *LBESignerT) AddAccounts(ctx *gin.Context, req *apistruct.AddAccountsReq) (*apistruct.AddAccountsResp, error) {
	var (
		resp = &apistruct.AddAccountsResp{}
	)
	accountSli, err := c.WalletCli.AddAccounts(ctx, req.BaseAddr, req.Chain, req.StartIndex, req.Count)
	if err != nil {
		log.ZError(ctx, "AddAccounts WalletCli.AddAccounts", err)
		return nil, err
	}
	resp.BaseAddr = req.BaseAddr
	for _, v := range accountSli {
		t := &apistruct.AccountDetail{}
		t.Addr = v.Addr
		t.Index = v.Index
		t.Chain = v.Chain
		resp.Adds = append(resp.Adds, t)
	}
	return resp, err
}

// AddDepositAddr add pooling addr
func (c *LBESignerT) AddDepositAddr(ctx *gin.Context, req *apistruct.AddDepositAddrReq) (*apistruct.AddDepositAddrResp, error) {
	var (
		resp = &apistruct.AddDepositAddrResp{}
	)
	err := c.WalletCli.AddDepositAddr(ctx, req.Addr, req.Name)
	if err != nil {
		log.ZError(ctx, "AddDepositAddr WalletCli.AddDepositAddr", err)
		return nil, err
	}
	resp.Addr = req.Addr
	return resp, err
}

// ListDepositAddr list pooling addr
func (c *LBESignerT) ListDepositAddr(ctx *gin.Context, req *apistruct.ListDepositAddrReq) (*apistruct.ListDepositAddrResp, error) {
	var (
		resp = &apistruct.ListDepositAddrResp{}
	)
	accounts, err := c.WalletCli.ListDepositAddr()
	if len(accounts) == 0 || err != nil {
		log.ZError(ctx, "AddDepositAddr WalletCli.AddDepositAddr", err)
		return nil, err
	}
	for _, v := range accounts {
		resp.DepositAccounts = append(resp.DepositAccounts, &apistruct.DepositAddrMeta{Addr: v.Addr, Name: v.Name})
	}
	return resp, err
}
