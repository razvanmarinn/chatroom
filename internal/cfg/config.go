package cfg

type Config struct {
	DbType    DatabaseType `yaml:"dbType"`
	CacheType CacheType    `yaml:"cacheType"`
	LogType   LogType      `yaml:"logType"`
}
type DatabaseType string
type CacheType string
type LogType string

const (
	Redis     CacheType = "redis"
	Memcached CacheType = "memcached"
)

const (
	PostgreSQL DatabaseType = "postgres"
)

const (
	Local       LogType = "local"
	Centralized LogType = "centralized"
)

// func LoadConfig() Config {
// 	var config Config

// 	file, err := os.Open("/cfg/config.yml")

// 	if err != nil {
// 		return config
// 	}
// 	defer file.Close()

// 	decoder := yaml.NewDecoder(file)
// 	err = decoder.Decode(&config)
// 	if err != nil {
// 		return config
// 	}

// 	return config
// }
func LoadConfig() Config {

	return Config{DbType: "postgres", CacheType: "redis", LogType: "local"}
}
