package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
)

//LambdaHandler . . .
func LambdaHandler(req *http.Request) error {
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
