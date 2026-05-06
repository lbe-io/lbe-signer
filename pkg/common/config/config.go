package config

type ApiConfig struct {
	Api struct {
		ListenIP string `yaml:"listenIP"`
		Ports    int    `yaml:"ports"`
		Swagger  bool   `yaml:"swagger"`
	} `yaml:"api"`
}

type Log struct {
	StorageLocation     string `yaml:"storageLocation"`
	RotationTime        uint   `yaml:"rotationTime"`
	RemainRotationCount uint   `yaml:"remainRotationCount"`
	RemainLogLevel      int    `yaml:"remainLogLevel"`
	IsStdout            bool   `yaml:"isStdout"`
	IsJson              bool   `yaml:"isJson"`
	IsSimplify          bool   `yaml:"isSimplify"`
	WithStack           bool   `yaml:"withStack"`
}

type Redis struct {
	ClusterMode bool     `yaml:"clusterMode"`
	Address     []string `yaml:"address"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	MaxRetry    int      `yaml:"maxRetry"`
	DB          int      `yaml:"db"`
	PoolSize    int      `yaml:"poolSize "` // Number of connections to pool.
}

type ZooKeeper struct {
	Schema   string   `yaml:"schema"`
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type Discovery struct {
	Enable    string    `yaml:"enable"`
	Etcd      Etcd      `yaml:"etcd"`
	ZooKeeper ZooKeeper `yaml:"zooKeeper"`
}

type Etcd struct {
	RootDirectory string   `yaml:"rootDirectory"`
	Address       []string `yaml:"address"`
	Username      string   `yaml:"username"`
	Password      string   `yaml:"password"`
}
