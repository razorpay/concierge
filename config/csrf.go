package config

//CSRF ...struct to hold CSRF configs
type CSRF struct {
	AuthKey string
	Secure  bool
}
