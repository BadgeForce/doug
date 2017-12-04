package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config struct that holds configs
type Config struct {
	AWSCONf   s3                       `toml:"s3"`
	Projects  []map[string]interface{} `toml:"projects"`
	Artifacts map[string]map[string]interface{}
}

//AWS . . .
type s3 struct {
	Regions []string `toml:"regions"`
	Bucket  string   `toml:"bucket"`
}

//Configs is the variable that holds config information for microservice
var Configs Config

func postDecode(config *Config) {
	config.Artifacts = make(map[string]map[string]interface{})
	for _, project := range config.Projects {
		name := project["name"].(string)
		config.Artifacts[name] = project
	}
}

//Initializes configurations
func init() {
	if _, err := toml.DecodeFile("./config.toml", &Configs); err != nil {
		log.Println("Check your configs.")
		log.Fatalln(err.Error())
		return
	}
	postDecode(&Configs)
	log.Print(Configs.Artifacts["badgeforce"])
}
