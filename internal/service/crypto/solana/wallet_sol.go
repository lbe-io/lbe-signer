package solana

import (
	"context"
	"fmt"
	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	bip32 "github.com/tyler-smith/go-bip32"
	"lbe_crypto_signer/internal/service/crypto/solana/hdwallets"
	"lbe_crypto_signer/pkg/tools/log"
	"strings"
)

type WalletSolana struct{}

func (w *WalletSolana) AccountsFromDeriveIndex(ctx context.Context, masterKey *bip32.Key, index int) (string, error) {
	return "", fmt.Errorf("derive from private key is not support")
}

func (w *WalletSolana) AccountsFromMnemonicIndex(ctx context.Context, Mnemonic string, index int) (string, error) {
	node := hdwallets.NewNode(Mnemonic, hdwallets.WithIndex(uint32(index)))
	return node.Address(), nil
}

func (w *WalletSolana) SignTx(ctx context.Context, Mnemonic string, addr string, index int, chainId string, needSign []byte) ([]byte, error) {
	node := hdwallets.NewNode(Mnemonic, hdwallets.WithIndex(uint32(index)))
	if strings.ToLower(node.Address()) != strings.ToLower(addr) {
		err := fmt.Errorf("node.Address == %s vs addr == %s", node.Address(), addr)
		log.ZError(ctx, "SignTx node.Address != addr ", err)
		return nil, err
	}

	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(needSign))
	if err != nil {
		log.ZError(ctx, "SignTx TransactionFromDecoder failed", err)
		return nil, err
	}
	txByte, err := tx.MarshalBinary()
	if err != nil {
		log.ZError(ctx, "SignTx tx.MarshalBinary failed", err)
		return nil, err
	}
	privateKey := solana.MustPrivateKeyFromBase58(node.SecretKey())
	signature, err := privateKey.Sign(txByte[:])
	if err != nil {
		log.ZError(ctx, "SignTx Sign failed", err)
		return nil, err
	}
	sigJsonByte, err := signature.MarshalJSON()
	if err != nil {
		log.ZError(ctx, "SignTx MarshalJSON failed", err)
		return nil, err
	}
	return sigJsonByte, nil
}
