package webhooklambda

import (
	"encoding/json"

	"github.com/BadgeForce/doug"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/github"
)

//LambdaHandler . . .
func lambdaHandler(req events.APIGatewayProxyRequest) error {
	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}
	}

	err = doug.UploadArtifacts(evt)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}
	}
	return return events.APIGatewayProxyResponse{
		Body: err.Error(),
		StatusCode: 400,
	}
}

func NewLamdaFn(configPath string) func(events.APIGatewayProxyRequest) error {
	doug.InitializeConfig(configPath)
	return lambdaHandler
}
