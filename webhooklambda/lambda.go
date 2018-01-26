package webhooklambda

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/BadgeForce/doug"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/github"
)

type S3UploadError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type ErrorResponseWrapper struct {
	Response events.APIGatewayProxyResponse
}

func (e *ErrorResponseWrapper) Error() string {
	b, _ := json.Marshal(e.Response)
	return string(b)
}

//LambdaHandler . . .
func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	hc, err := doug.ParseHook([]byte(doug.Configs.Github.Secret), req.Headers, req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	errs := doug.UploadArtifacts(evt)
	if errs != nil {
		return events.APIGatewayProxyResponse{}, errors.New(fmt.Sprintf("%+v", errs))
	}

	return getGateWayRes("Artifacts Uploaded", 200), nil
}

func getGateWayRes(body string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: statusCode,
	}
}

func NewLamdaFn(configPath string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	doug.InitializeConfig(configPath)
	return lambdaHandler
}
