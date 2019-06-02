package utils

const (
	CONFIGPATH = "/conf/config.toml"
)

type DefaultTomlConfig struct {
	Server struct {
		Address      string `toml:"address"`
		Port         int    `toml:"port"`
		OpsPort      int    `toml:"ops_port"`
		ServiceName  string `toml:"service_name"`
		Degrading    bool   `toml:"degrading"`
		DegradingAPI string `toml:"degrading_api"`
	} `toml:"server"`
	Log struct {
		Level string `toml:"level"`
		Path  string `toml:"path"`
	} `toml:"log"`
	ServerClient []struct {
		EndPoint    string `toml:"end_point"`
		ServiceName string `toml:"service_name"`
		Port        int    `toml:"port"`
		OpsPort        int    `toml:"ops_port"`
	} `toml:"server-client"`
	Mysql struct {
		URL      string `toml:"url"`
		PoolSize int    `toml:"pool_size"`
	} `toml:"mysql"`
	Redis []struct {
		EndPoint  string `toml:"end_point"`
		Port      int    `toml:"port"`
		CacheName string `toml:"cache_name"`
	} `toml:"redis"`
}
