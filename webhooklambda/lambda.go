package webhooklambda

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/BadgeForce/doug"
)

//LambdaHandler . . .
func lambdaHandler(req *http.Request) error {
	hc, err := ParseHook([]byte(doug.Configs.Github.Secret), req)
	if err != nil {
		return err
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return err
	}

	return nil
}

func GetLamdaFn(configPath string) fn {
	doug.InitializeConfig(configPath)
	return lambdaHandler
}
