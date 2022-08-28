package config

// RedisConfig ....
type RedisConfig struct {
	Host      string
	Port      string
	Password  string
	Database  int
	MaxIdle   int
	MaxActive int
}
