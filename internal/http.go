package internal

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var client = resty.New()
var clientVersion = "dev"

func init() {
	if base, found := os.LookupEnv("CNATIVE_API"); found {
		client.SetBaseURL(base)
	} else {
		client.SetBaseURL("https://api.cnative.dev")
	}
}

type ErrResp struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func SetVersion(version string) {
	clientVersion = version
}

func R() *resty.Request {
	request := client.R()
	request.SetHeader("Content-Type", "application/json")
	request.SetHeader("User-Agent", fmt.Sprintf("cnative/%s", clientVersion))
	request.SetError(ErrResp{})
	if viper.IsSet("token") {
		tokenString := viper.GetString("token")
		if len(tokenString) > 0 {
			if token, _ := jwt.Parse(tokenString, nil); token != nil && token.Claims.Valid() == nil {
				request.SetAuthToken(tokenString)
			} else {
				viper.Set("token", "")
			}
		}
	}
	return request
}
