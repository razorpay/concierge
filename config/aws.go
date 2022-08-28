package config

type SecurityGroupIngress struct {
	RuleType string
	PortTo   int64
	PortFrom int64
	Protocol string
}

type S3BucketPolicy struct {
	Version   string
	Id        string `json:",omitempty"`
	Statement []struct {
		Sid          string
		Effect       string
		NotPrincipal interface{} `json:",omitempty"`
		Principal    interface{} `json:",omitempty"`
		Action       string
		Resource     interface{} `json:",omitempty"`
		Condition    map[string]interface{}
	}
}
