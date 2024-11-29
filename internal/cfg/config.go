package cfg

type Config struct {
	DbType DatabaseType
	CacheType CacheType 
}
type DatabaseType string
type CacheType string

const (
	Redis CacheType = "redis"
	Memcached CacheType = "memcached"

)

const (
	PostgreSQL DatabaseType = "postgres"
)

// TODO: Implement config
func LoadConfig() Config {
	return Config{DbType: "postgres", CacheType: "redis"}
}
