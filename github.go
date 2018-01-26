package doug

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
)

const (
	sigHeader        = "X-Hub-Signature"
	ghEventHeader    = "X-GitHub-Event"
	ghDeliveryHeader = "X-GitHub-Delivery"
	filePathPrefix   = "/tmp/"
)

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
func ParseHook(secret []byte, requestHeaders map[string]string, requestBody string) (*HookContext, error) {
	hc := HookContext{}
	if hc.Signature = requestHeaders[sigHeader]; len(hc.Signature) == 0 {
		return nil, errors.New("No signature!")
	}

	if hc.Event = requestHeaders[ghEventHeader]; len(hc.Event) == 0 {
		return nil, errors.New("No event!")
	}

	if hc.Id = requestHeaders[ghDeliveryHeader]; len(hc.Id) == 0 {
		return nil, errors.New("No event Id!")
	}

	body := []byte(requestBody)

	if !verifySignature(secret, hc.Signature, body) {
		return nil, errors.New("Invalid signature")
	}

	hc.Payload = body

	return &hc, nil
}

func makeTempDir() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	tmpDir := filePathPrefix + id.String()

	err = os.Mkdir(tmpDir, 0777)
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func removeTempDir(name string) {
	os.RemoveAll(name)
}

//CloneRepo . . .
func CloneRepo(url string, tagName string) (string, error) {
	dir, err := makeTempDir()
	if err != nil {
		return "", err
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return "", err
	}

	return dir, nil
}
