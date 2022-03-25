package config

type SecurityGroupIngress struct {
	RuleType string
	PortTo   int64
	PortFrom int64
	Protocol string
}

type S3BucketPolicy struct {
	Version   string
	ID        string `json:",omitempty"`
	Statement []struct {
		Sid          string
		Effect       string
		NotPrincipal map[string]interface{} `json:",omitempty"`
		Principal    map[string]interface{} `json:",omitempty"`
		Action       string
		Resource     []string
		Condition    map[string]interface{}
	}
}
