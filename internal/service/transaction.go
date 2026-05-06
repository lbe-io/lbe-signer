package service

import (
	"github.com/gin-gonic/gin"
	"lbe_crypto_signer/pkg/common/apistruct"
	"lbe_crypto_signer/pkg/tools/log"
)

// Transaction Sign
func (c *LBESignerT) Sign(ctx *gin.Context, req *apistruct.SignerTxReq) (*apistruct.SignerTxResp, error) {
	signRes, err := c.WalletCli.Sign(ctx, req.BaseAddr, req.Chain, req.SignAddr, req.ChainID, req.AddrIndex, req.SignData)
	if err != nil {
		log.ZError(ctx, "c.WalletCli.Sign", err)
		return nil, err
	}
	return &apistruct.SignerTxResp{SignTx: signRes}, err
}
