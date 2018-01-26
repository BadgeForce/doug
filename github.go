package doug

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	uuid "github.com/satori/go.uuid"
)

const branch = "--branch"
const depth = "--depth"

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

//HookContext . . .
type HookContext struct {
	Signature string
	Event     string
	Id        string
	Payload   []byte
}

//ParseHook . . .
func ParseHook(secret []byte, request events.APIGatewayProxyRequest) (*HookContext, error) {
	// hc := HookContext{}
	fmt.Printf("%+v\n", request.Headers)
	fmt.Printf("%+v\n", request.Body)
	return nil, errors.New("messed up")
	// if hc.Signature = req.Header.Get("x-hub-signature"); len(hc.Signature) == 0 {
	// 	return nil, errors.New("No signature!")
	// }
	//
	// if hc.Event = req.Header.Get("x-github-event"); len(hc.Event) == 0 {
	// 	return nil, errors.New("No event!")
	// }
	//
	// if hc.Id = req.Header.Get("x-github-delivery"); len(hc.Id) == 0 {
	// 	return nil, errors.New("No event Id!")
	// }
	//
	// body, err := ioutil.ReadAll(req.Body)
	//
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !verifySignature(secret, hc.Signature, body) {
	// 	return nil, errors.New("Invalid signature")
	// }
	//
	// hc.Payload = body
	//
	// return &hc, nil
}

func makeTempDir() (string, error) {
	name, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	err = os.Mkdir(name.String(), 0777)
	if err != nil {
		return "", err
	}

	return name.String(), nil
}

func removeTempDir(name string) {
	os.RemoveAll("./" + name)
}

//CloneRepo . . .
func CloneRepo(sshURL string, tagName string) (string, error) {
	dir, err := makeTempDir()
	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", "clone", branch, tagName, depth, "1", sshURL, dir)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return dir, nil
}
