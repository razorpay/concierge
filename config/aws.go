package config

type SecurityGroupIngress struct {
	RuleType string
	PortTo   int64
	PortFrom int64
	Protocol string
}
