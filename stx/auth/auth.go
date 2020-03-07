package auth

// Auth holds fields related to aws-vault and ykman config
type Auth struct {
	AwsVault struct {
		Enabled bool
	}
	Ykman struct {
		Enabled bool
	}
}
