package config

//Application ...struct to hold application level configs
type Application struct {
	Name           string
	Mode           string
	ListenPort     string
	ListenIP       string
	MaxExpiry      int
	CookieSecure   bool
	CookieHTTPOnly bool
}
