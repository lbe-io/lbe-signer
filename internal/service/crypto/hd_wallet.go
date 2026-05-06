package crypto

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"go.uber.org/zap"
	"lbe_crypto_signer/internal/service/crypto/ethereum"
	"lbe_crypto_signer/internal/service/crypto/solana"
	"lbe_crypto_signer/internal/service/crypto/tron"
	"lbe_crypto_signer/pkg/tools/log"
	"strings"
	"sync"
	"time"
)

const (
	ETH_TYP_ACCOUNT    = "ethereum"
	BSC_TYP_ACCOUNT    = "bsc"
	TRON_TYP_ACCOUNT   = "tron"
	SOLANA_TYP_ACCOUNT = "solana"
)

const MasterAccountIndex = 0

type WalletBehaviors interface {
	AccountsFromDeriveIndex(ctx context.Context, masterKey *bip32.Key, index int) (string, error)
	AccountsFromMnemonicIndex(ctx context.Context, Mnemonic string, index int) (string, error)
	SignTx(ctx context.Context, Mnemonic string, addr string, index int, chainId string, needSign []byte) ([]byte, error)
}

type AccountDetail struct {
	Addr  string
	Chain string
	Index int64
}

type AccountsMeta struct {
	Key          string
	AccountsComp *AccountsComp
}

type AccountsComp struct {
	Name     string
	Accounts []AccountsBase
}

type AccountsBase struct {
	Addr  string `json:"addr"`
	Chain string `json:"chain"`
}

type DepositAddrMeta struct {
	Name string
	Addr string
}

type WalletBase struct {
	KeyRamStore      sync.Map
	KeyRelationStore sync.Map
	DepositAddrStore sync.Map
	EthWallet        WalletBehaviors
	TronWallet       WalletBehaviors
	SolanaWallet     WalletBehaviors
}

func NewWalletBase() *WalletBase {
	return &WalletBase{
		EthWallet:    &ethereum.WalletEthereum{},
		TronWallet:   &tron.WalletTron{},
		SolanaWallet: &solana.WalletSolana{},
	}
}

func (w *WalletBase) GenerateMnemonic(ctx context.Context, size MnemonicBitSize) (string, error) {
	entropy, err := bip39.NewEntropy(int(size))
	if err != nil {
		log.ZError(ctx, "GenerateMnemonic NewEntropy failed:", err)
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.ZError(ctx, "GenerateMnemonic NewMnemonic failed:", err)
		return "", err
	}
	return mnemonic, nil
}

func (w *WalletBase) SeedToMasterKey(ctx context.Context, seed string, passPhase string) (*bip32.Key, error) {
	keySeeds := bip39.NewSeed(seed, passPhase)
	masterKey, err := bip32.NewMasterKey(keySeeds)
	if err != nil {
		log.ZError(ctx, "SeedToMastKey NewMasterKey failed:", err)
		return nil, err
	}
	return masterKey, nil
}

func (w *WalletBase) AccountsFromDerive(ctx context.Context, masterKey *bip32.Key, Mnemonic string, name string) ([]AccountsBase, error) {
	ethAddr, err := w.EthWallet.AccountsFromDeriveIndex(ctx, masterKey, MasterAccountIndex)
	if err != nil {
		log.ZError(ctx, "EthWallet AccountsFromDeriveIndex", err)
		return nil, err
	}
	tronAddr, _ := w.TronWallet.AccountsFromDeriveIndex(ctx, masterKey, MasterAccountIndex)
	if err != nil {
		log.ZError(ctx, "TronWallet AccountsFromDeriveIndex", err)
		return nil, err
	}
	solAddr, _ := w.SolanaWallet.AccountsFromMnemonicIndex(ctx, Mnemonic, MasterAccountIndex)
	if err != nil {
		log.ZError(ctx, "SolanaWallet AccountsFromDeriveIndex", err)
		return nil, err
	}
	log.ZInfo(ctx, "AccountsFromDerive", zap.Any("ethAddr", ethAddr), zap.Any("tronAddr", tronAddr), zap.Any("solAddr", solAddr))
	accountSli := []AccountsBase{{Addr: ethAddr, Chain: ETH_TYP_ACCOUNT}, {Addr: tronAddr, Chain: TRON_TYP_ACCOUNT}, {Addr: solAddr, Chain: SOLANA_TYP_ACCOUNT}}
	w.KeyRamStore.Store(ethAddr, Mnemonic)
	w.KeyRamStore.Store(tronAddr, Mnemonic)
	w.KeyRamStore.Store(solAddr, Mnemonic)
	MnemonicSeed := fmt.Sprintf("%s%d", Mnemonic, time.Now().UnixMilli())
	md := md5.Sum([]byte(MnemonicSeed))
	mdStr := hex.EncodeToString(md[:])
	comp := &AccountsComp{
		Name:     name,
		Accounts: accountSli,
	}
	w.KeyRelationStore.Store(mdStr, comp)

	return accountSli, nil
}

func (w *WalletBase) KeyList() ([]AccountsMeta, error) {
	var metaSli []AccountsMeta
	w.KeyRelationStore.Range(func(key, value any) bool {
		a := AccountsMeta{}
		a.Key = key.(string)
		a.AccountsComp = value.(*AccountsComp)
		metaSli = append(metaSli, a)
		return true
	})
	return metaSli, nil
}

func (w *WalletBase) AddAccounts(ctx context.Context, key string, chain string, start int64, count int64) ([]*AccountDetail, error) {
	Mnemonic, exist := w.KeyRamStore.Load(key)
	if !exist {
		err := fmt.Errorf("Key : %s is not exist.", key)
		log.ZError(ctx, "KeyRamStore Load failed", err)
		return nil, err
	}

	handler, err := w.selectHandler(chain)
	if err != nil {
		return nil, err
	}
	MnemonicStr, ok := Mnemonic.(string)
	if !ok {
		err = fmt.Errorf("Mnemonic is invalid format.")
		log.ZError(ctx, "KeyRamStore Load failed", err)
		return nil, err
	}
	var AccountSli []*AccountDetail
	for i := start; i < start+count; i++ {
		addr, er := handler.AccountsFromMnemonicIndex(ctx, MnemonicStr, int(i))
		if er != nil {
			return nil, er
		}
		t := &AccountDetail{
			Addr:  addr,
			Index: i,
			Chain: chain,
		}
		AccountSli = append(AccountSli, t)
	}

	return AccountSli, nil
}

func (w *WalletBase) AddDepositAddr(ctx context.Context, addr string, name string) error {
	exist := false
	w.DepositAddrStore.Range(func(key, value any) bool {
		keyStr, _ := key.(string)
		if strings.ToLower(keyStr) == strings.ToLower(addr) {
			exist = true
			return false
		}
		return true
	})
	if exist {
		err := fmt.Errorf("Key : %s is exist.", addr)
		log.ZError(ctx, "DepositAddrStore ", err)
		return err
	}

	w.DepositAddrStore.Store(addr, name)
	return nil
}

func (w *WalletBase) ListDepositAddr() ([]*DepositAddrMeta, error) {
	var depositSlice []*DepositAddrMeta
	w.DepositAddrStore.Range(func(key, value any) bool {
		keyStr, _ := key.(string)
		valueStr, _ := value.(string)
		depositSlice = append(depositSlice, &DepositAddrMeta{Addr: keyStr, Name: valueStr})
		return true
	})
	return depositSlice, nil
}

func (w *WalletBase) Sign(ctx context.Context, key string, chain string, chainId string, addr string, index int64, needSign []byte) ([]byte, error) {
	Mnemonic, exist := w.KeyRamStore.Load(key)
	if !exist {
		err := fmt.Errorf("Sign Key : %s is not exist.", key)
		log.ZError(ctx, "KeyRamStore Load failed", err)
		return nil, err
	}

	handler, err := w.selectHandler(chain)
	if err != nil {
		return nil, err
	}
	MnemonicStr, ok := Mnemonic.(string)
	if !ok {
		err = fmt.Errorf("Mnemonic is invalid format.")
		log.ZError(ctx, "KeyRamStore Load failed", err)
		return nil, err
	}

	signRes, er := handler.SignTx(ctx, MnemonicStr, addr, int(index), chainId, needSign)
	if er != nil {
		return nil, er
	}

	return signRes, nil
}

func (w *WalletBase) selectHandler(chain string) (WalletBehaviors, error) {
	if chain == ETH_TYP_ACCOUNT {
		return w.EthWallet, nil
	} else if chain == TRON_TYP_ACCOUNT {
		return w.TronWallet, nil
	} else if chain == SOLANA_TYP_ACCOUNT {
		return w.SolanaWallet, nil
	}
	return nil, fmt.Errorf("Chain: %s is not support", chain)
}
