package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/github"
)

func Handler(req *http.Request) error {
	hc, err := ParseHook([]byte(Configs.Github.Secret), req)
	if err != nil {
		return err
	}

	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(hc.Payload, &evt); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
