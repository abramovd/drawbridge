package jwt

import (
	"os"
	"fmt"
	"errors"
	"net/http"

	"github.com/jakewright/drawbridge/plugin"
	"github.com/jakewright/muxinator"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

func init() {
	plugin.RegisterPlugin("jwt", &JWTAuthorizer{})
}

type JWTAuthorizer struct{}

// Middleware to check JWT authorization on request
func (r *JWTAuthorizer) Middleware(cfg map[string]interface{}) (muxinator.Middleware, error) {
	var opts Options
	if err := plugin.DecodeConfig(cfg, &opts); err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		issSecretMap := make(map[string][]byte)
		for _, issuer := range opts.Issuers {
			issSecretMap[issuer.Iss] = []byte(os.Getenv(issuer.SecretEnvPath))
		}

		token, err := request.ParseFromRequest(
			r, 
			request.AuthorizationHeaderExtractor, 
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					issuer := claims["iss"].(string)
					if secret, ok := issSecretMap[issuer]; !ok {
						return secret, nil
					}
				}

				return nil, errors.New("Unsupported iss")
		})
	
		if err == nil && token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}, nil
}