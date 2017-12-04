package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/go-github/github"
)

const truffleArtifactPath = "/build/contracts"

//UploadArtifacts . . .
func UploadArtifacts(event github.ReleaseEvent) error {
	repo := event.GetRepo()
	tagName := event.GetRelease().GetTagName()
	sshURL := repo.GetSSHURL()
	dir, err := CloneRepo(sshURL, tagName)
	if err != nil {
		return err
	}
	log.Println(dir)
	log.Println(Configs)

	//upload(dir, repo.GetName()+"/"+tagName+)
	return nil
}

func upload(tempDir string, s3Path string) error {
	file, err := os.Open("./" + tempDir + truffleArtifactPath + "/Issuer.json")
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	go func() {
		content := bufio.NewWriter(writer)
		io.Copy(content, file)
		file.Close()
		writer.Close()
	}()

	creds := credentials.NewEnvCredentials()
	for _, region := range Configs.AWSCONf.Regions {
		go func(creds *credentials.Credentials, region string, bucket string) {
			uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(region), Credentials: creds}))
			_, err = uploader.Upload(&s3manager.UploadInput{
				Body:   reader,
				Bucket: aws.String(bucket),
				Key:    aws.String(string(s3Path)),
			})
			if err != nil {
				log.Println(err.Error())
			}
		}(creds, region, Configs.AWSCONf.Bucket)
	}

	return nil
}
