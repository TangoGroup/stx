package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/logrusorgru/aurora"
)

// vault is used to retrieve credentials from aws-vault
type vault struct {
	awsConfig  *ConfigFile
	authConfig Auth
	cache      map[string]*credentials.Credentials
}

func newVault(awsConfig *ConfigFile, authConfig Auth) *vault {
	v := vault{
		awsConfig:  awsConfig,
		authConfig: authConfig,
		cache:      make(map[string]*credentials.Credentials),
	}
	return &v
}

func (v *vault) getCredentials(profile string) *credentials.Credentials {
	if credentials, ok := v.cache[profile]; ok {
		return credentials
	}

	// TODO move this to log
	au := aurora.NewAurora(true)

	_, existingVault := os.LookupEnv("AWS_VAULT")
	if existingVault {
		// TODO move to log.Fatal
		fmt.Println(au.Red("Cannot run in nested aws-vault session!"))
		os.Exit(1)
	}

	profileSection, profileExists := v.awsConfig.ProfileSection(profile)
	if !profileExists {
		fmt.Println(au.Red(profile + " does not exist in .aws/config"))
		os.Exit(1)
	}

	//fmt.Printf("%+v\n", profileSection)

	execArgs := []string{"exec", "--no-session"}

	if profileSection.MfaSerial != "" {

		var mfa string

		if v.authConfig.Ykman.Enabled {
			sourceProfile := profile
			if profileSection.SourceProfile != "" {
				sourceProfile = profileSection.SourceProfile
			}

			ykmanOutput, ykmanErr := exec.Command("ykman", "oath", "code", "-s", sourceProfile).Output()
			if ykmanErr != nil {
				fmt.Println(au.Red(ykmanErr))
				os.Exit(1)
			}
			mfa = strings.TrimSpace(string(ykmanOutput))
			// TODO move to log.Debug
			fmt.Println("Pulled MFA from ykman profile", sourceProfile)
		} else {
			fmt.Print("MFA: ")
			fmt.Scanln(&mfa)
		}
		execArgs = append(execArgs, "-t", mfa)
	}

	execArgs = append(execArgs, "--json", profile)
	// get credentials from aws-vault
	execOut, execErr := exec.Command("aws-vault", execArgs...).Output()
	if execErr != nil {
		fmt.Println(au.Red(execErr))
		os.Exit(1)
	}

	// TODO: cache credentials until expired
	var rawCreds map[string]interface{}
	json.Unmarshal(execOut, &rawCreds)
	creds := credentials.NewStaticCredentials(fmt.Sprint(rawCreds["AccessKeyId"]), fmt.Sprint(rawCreds["SecretAccessKey"]), fmt.Sprint(rawCreds["SessionToken"]))
	return creds
}
