package doug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/github"
)

func getTestEvnt() github.ReleaseEvent {
	file, e := ioutil.ReadFile("./test-files/gh-release-payload.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	evt := github.ReleaseEvent{}
	if err := json.Unmarshal(file, &evt); err != nil {
		log.Fatalln("Failed processing test hook! ('%s')", err)
	}

	return evt
}
