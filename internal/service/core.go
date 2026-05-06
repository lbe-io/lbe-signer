package service

import (
	"context"
	"lbe_crypto_signer/internal/service/crypto"
	"lbe_crypto_signer/pkg/common/config"
)

type ActiveConfig struct {
	ApiConf *config.ApiConfig
}

type LBESignerT struct {
	Conf      *ActiveConfig
	WalletCli *crypto.WalletBase
}

func NewLBESignerT(conf *ActiveConfig) (*LBESignerT, error) {
	t := &LBESignerT{
		Conf: conf,
	}
	err := t.initWalletCli(context.Background())
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (l *LBESignerT) initWalletCli(ctx context.Context) error {
	l.WalletCli = crypto.NewWalletBase()
	return nil
}

var Core *LBESignerT

func Init(config *ActiveConfig) {
	var err error
	Core, err = NewLBESignerT(config)
	if err != nil {
		panic(err)
	}
}
