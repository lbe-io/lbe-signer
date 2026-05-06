package apistruct

type SignerTxReq struct {
	Chain     string `json:"chain"`
	ChainID   string `json:"chain_id"`
	BaseAddr  string `json:"base_addr"`
	SignAddr  string `json:"sign_addr"`
	AddrIndex int64  `json:"addr_index"`
	SignData  []byte `json:"sign_data"`
}

type SignerTxResp struct {
	SignTx []byte `json:"sign_tx"`
}
