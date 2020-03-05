package auth

import "github.com/aws/aws-sdk-go/aws/credentials"

// credentialStore is used to cache credentials with the profile they were pulled from
type credentialStore struct {
	cache map[string]*credentials.Credentials
}

func newCredentialStore() *credentialStore {
	cc := credentialStore{}
	return &cc
}

func (cc *credentialStore) getCredentials(profile string) *credentials.Credentials {
	if credentials, ok := cc.cache[profile]; ok {
		return credentials
	}
	return nil
}

func (cc *credentialStore) setCredentials(profile string, credentials *credentials.Credentials) {
	cc.cache[profile] = credentials
}
