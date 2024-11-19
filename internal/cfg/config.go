package cfg

type Config struct {
	DbType DatabaseType
}
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres"
)

// TODO: Implement config
func LoadConfig() Config {
	return Config{DbType: "postgres"}
}
