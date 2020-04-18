package jwt

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/jakewright/drawbridge/utils"
)

var supportedSigningAlgorithms = []string {
	jwt.SigningMethodHS256.Alg(), 
	jwt.SigningMethodHS384.Alg(), 
	jwt.SigningMethodHS512.Alg(),
}


type issuerOptions struct {
	Iss           string    `mapstructure:"iss"`
	SecretEnvPath string `mapstructure:"secret_env_path"`
}

// Options holds config for the jwt authorizer middleware
type Options struct {
	Alrorithm string   `mapstructure:"algorithm"`
	Issuers   []issuerOptions `mapstructure:"issuers"`
}

// Validate returns an error if there are any problems with the config
func (o *Options) Validate() error {
	if o == nil {
		return errors.New("config is nil")
	}

	if len(o.Issuers) == 0 {
		return errors.New("At least 1 issuer config must be set")
	}

	if !utils.StringInSlice(o.Alrorithm, supportedSigningAlgorithms) {
		return fmt.Errorf(
			"%s signing algorithm is not supported. Supported algorithms: %s", 
			o.Alrorithm, supportedSigningAlgorithms,
		)
	}

	issuersMap := make(map[string]bool)
	for _, issuer := range o.Issuers {

		// Check iss
		if issuer.Iss == "" {
			return errors.New("iss must be set")
		}
		if issuersMap[issuer.Iss] {
			return fmt.Errorf("iss must be unique. iss %s is duplicate", issuer.Iss)
		}

		// check secret can be found
		secret := os.Getenv(issuer.SecretEnvPath)
		if secret == "" {
			return fmt.Errorf("Environment variable %s for secret not found", issuer.SecretEnvPath)
		}
	}

	return nil
}
