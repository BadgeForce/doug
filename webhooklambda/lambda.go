package webhooklambda

import (
	"encoding/json"
	"net/http"

	"github.com/BadgeForce/doug"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/github"
)

//LambdaHandler . . .
func lambdaHandler(request events.APIGatewayProxyRequest) error {
	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), req)
	if err != nil {
		return err
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return err
	}

	return nil
}

func NewLamdaFn(configPath string) func(*http.Request) error {
	doug.InitializeConfig(configPath)
	return lambdaHandler
}
