package config

//DatabaseConfig ...
type DatabaseConfig struct {
	Host         string
	DBName       string
	DBUsername   string
	DBPassword   string
	DBPort       string
	DBDatabase   string
	MaxIdleConns int
	MaxOpenConns int
}
