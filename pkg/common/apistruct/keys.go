package apistruct

type KeysAddMnemonicReq struct {
	Mnemonic string `json:"mnemonic"`
	Name     string `json:"name"`
}

type KeysAddMnemonicResp struct {
	Accounts []AccountBasic `json:"accounts"`
}

type AccountBasic struct {
	Addr  string `json:"addr"`
	Chain string `json:"chain"`
}

type KeysListReq struct{}

type KeysListResp struct {
	KeySlice []KeyMeta `json:"keySlice"`
}

type KeyMeta struct {
	KeyFLag  string         `json:"keyFLag"`
	Name     string         `json:"name"`
	Accounts []AccountBasic `json:"accounts"`
}

type AddAccountsReq struct {
	BaseAddr   string `json:"baseAddr"`
	Chain      string `json:"chain"`
	StartIndex int64  `json:"startIndex"`
	Count      int64  `json:"count"`
}

type AddAccountsResp struct {
	BaseAddr string           `json:"baseAddr"`
	Adds     []*AccountDetail `json:"adds"`
}

type AddDepositAddrReq struct {
	Addr string `json:"addr"`
	Name string `json:"name"`
}

type AddDepositAddrResp struct {
	Addr string `json:"baseAddr"`
}

type ListDepositAddrReq struct{}

type ListDepositAddrResp struct {
	DepositAccounts []*DepositAddrMeta `json:"depositAccounts"`
}

type DepositAddrMeta struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type AccountDetail struct {
	Addr  string `json:"addr"`
	Index int64  `json:"index"`
	Chain string `json:"chain"`
}
