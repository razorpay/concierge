package config

type Looker struct {
	ClientId      string
	ClientSecret  string
	BaseUrl       string
	IsEnabled     bool
	DatumHostname string
	DatumAuthSecret string
}
