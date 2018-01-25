package doug

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

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
	upload(dir, repo.GetName()+"/"+tagName, repo.GetName())
	return nil
}

func upload(tempDir string, s3Path string, project string) {
	var wg sync.WaitGroup
	for _, artifact := range Configs.Artifacts[project] {
		path := fmt.Sprintf("%s/%s", s3Path, artifact.(string))
		err := s3Upload(tempDir, path, artifact.(string), &wg)
		if err != nil {
			log.Println(err.Error())
		}
	}
	wg.Wait()
	removeTempDir(tempDir)
}

func s3Upload(tempDir string, s3Path string, artifact string, wg *sync.WaitGroup) error {
	file, err := os.Open("./" + tempDir + truffleArtifactPath + "/" + artifact)
	if err != nil {
		return err
	}
	reader, writer := io.Pipe()
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		content := bufio.NewWriter(writer)
		io.Copy(content, file)
		file.Close()
		writer.Close()
	}(wg)

	putToRegions(reader, s3Path, wg)

	return nil
}

// loop through configured s3 regions and put artifact in each region
func putToRegions(reader *io.PipeReader, s3Path string, wg *sync.WaitGroup) {
	creds := credentials.NewEnvCredentials()
	for _, region := range Configs.S3Conf.Regions {
		wg.Add(1)
		go func(creds *credentials.Credentials, region string, bucket string, wg *sync.WaitGroup) {
			defer wg.Done()
			putObj(creds, region, bucket, s3Path, reader)
		}(creds, region, Configs.S3Conf.Bucket, wg)
	}
}

// put artifact to specific region
func putObj(creds *credentials.Credentials, region string, bucket string, s3Path string, reader *io.PipeReader) {
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(region), Credentials: creds}))
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(bucket),
		Key:    aws.String(string(s3Path)),
	})
	if err != nil {
		log.Println(err.Error())
	}
}
