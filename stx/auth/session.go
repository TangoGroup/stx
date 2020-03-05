package auth

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionProvider provides a single reference to retrieve sessions. It retrieves and caches credentials privately.
type SessionProvider struct {
	awsConfig  *ConfigFile
	authConfig Auth
	vault      *vault
	cache      map[string]*session.Session
}

// NewSessionProvider provides a way for commands to retrieve sessions based on stack.Profile
func NewSessionProvider(authConfig Auth) *SessionProvider {

	sp := SessionProvider{
		authConfig: authConfig,
		cache:      make(map[string]*session.Session),
	}

	if authConfig.AwsVault.Enabled {
		// load ~/.aws/config
		awsConfig, awsConfigErr := LoadConfigFromEnv()
		if awsConfigErr != nil {
			fmt.Println(awsConfigErr)
			os.Exit(1)
		}
		sp.awsConfig = awsConfig
		sp.vault = newVault(awsConfig, authConfig)
	}

	return &sp
}

// GetSession returns aws session using credentials associated with stack.Profile
func (sp *SessionProvider) GetSession(profile string) (*session.Session, error) {

	if sess, ok := sp.cache[profile]; ok {
		return sess, nil
	}

	var sess *session.Session
	var sessErr error

	if sp.authConfig.AwsVault.Enabled {
		creds := sp.vault.getCredentials(profile)
		config := aws.NewConfig().WithCredentials(creds)
		sess, sessErr = session.NewSession(config)
	} else {
		// TODO: test this with keys in ~/.aws/credentials
		sess, sessErr = session.NewSessionWithOptions(session.Options{Profile: profile})
	}

	sp.cache[profile] = sess
	return sess, sessErr
}
