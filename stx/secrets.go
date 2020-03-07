package stx

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"go.mozilla.org/sops/v3/decrypt"
)

// DecryptSecrets uses sops to decrypt the file with credentials from the given profile
func DecryptSecrets(file string, creds *credentials.Credentials) (map[string]string, error) {
	// set ENV vars (primarily for sops decrypt)
	value, _ := creds.Get()
	os.Setenv("AWS_ACCESS_KEY_ID", value.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", value.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", value.SessionToken)
	sopsOutput, sopsError := decrypt.File(file, "Dotenv")

	secrets, _ := godotenv.Parse(bytes.NewReader(sopsOutput))

	return secrets, sopsError
}
