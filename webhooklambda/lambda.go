package webhooklambda

import (
	"encoding/json"
	"fmt"

	"github.com/BadgeForce/doug"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/github"
)

type S3UploadError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

//LambdaHandler . . .
func lambdaHandler(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), req.Headers, req.Body)
	if err != nil {
		return getGateWayRes(err.Error(), 400)
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return getGateWayRes(err.Error(), 400)
	}

	errors := doug.UploadArtifacts(evt)
	if errors != nil {
		b, _ := json.Marshal(S3UploadError{
			"Errors while uploading artifacts to S3",
			[]string{fmt.Sprintf("%+v", errors)},
		})
		return getGateWayRes(string(b), 500)
	}

	return getGateWayRes("Artifacts Uploaded", 200)
}

func getGateWayRes(body string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: statusCode,
	}
}

func NewLamdaFn(configPath string) func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	doug.InitializeConfig(configPath)
	return lambdaHandler
}
