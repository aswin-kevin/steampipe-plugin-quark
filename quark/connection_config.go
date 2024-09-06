package quark

type awsConfig struct {
	AwsAccessKeyID     string `hcl:"access_key"`
	AwsSecretAccessKey string `hcl:"secret_key"`
	AwsSessionToken    string `hcl:"session_token"`
}

func ConfigInstance() interface{} {
	return &awsConfig{}
}
