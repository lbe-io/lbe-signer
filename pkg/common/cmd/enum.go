package cmd

const (
	FlagConf = "config-path"
)

var (
	RedisConfigFileName   string
	MongodbConfigFileName string
	LogConfigFileName     string
	ScanCfgFileName       string
	ApiCfgFileName        string
	ChainCfgFileName      string
	CoinCfgFileName       string
)

func init() {
	RedisConfigFileName = "redis.yaml"
	MongodbConfigFileName = "mongodb.yaml"
	LogConfigFileName = "log.yaml"
	ScanCfgFileName = "active.yaml"
	ApiCfgFileName = "api.yaml"
	ChainCfgFileName = "chain.yaml"
	CoinCfgFileName = "coin.yaml"
}
